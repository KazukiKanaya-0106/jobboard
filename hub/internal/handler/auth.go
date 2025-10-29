package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
	"github.com/kanaya/jobboard-hub/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	queries   repo.Querier
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewAuthHandler(queries repo.Querier, jwtSecret []byte, tokenTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		queries:   queries,
		jwtSecret: jwtSecret,
		tokenTTL:  tokenTTL,
	}
}

type authRequest struct {
	ClusterId string `json:"cluster_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type authResponse struct {
	ClusterId string `json:"cluster_id"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	cluster, err := h.queries.GetCluster(c.Request.Context(), req.ClusterId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cluster.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	resp, err := h.issueTokenResponse(req.ClusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	_, err = h.queries.CreateCluster(c.Request.Context(), repo.CreateClusterParams{
		ID:           req.ClusterId,
		PasswordHash: string(hashed),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			c.JSON(http.StatusConflict, gin.H{"error": "cluster already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cluster"})
		return
	}

	resp, err := h.issueTokenResponse(req.ClusterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) issueTokenResponse(clusterId string) (authResponse, error) {
	now := time.Now()
	claims := middleware.AuthClaims{
		ClusterId: clusterId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(h.tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return authResponse{}, err
	}
	return authResponse{
		ClusterId: clusterId,
		Token:     signed,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

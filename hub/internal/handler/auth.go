package handler

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kanaya/jobboard-hub/internal/apierror"
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
	ClusterID string `json:"cluster_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type authResponse struct {
	ClusterID string `json:"cluster_id"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Write(c, apierror.InvalidRequest)
		return
	}

	cluster, err := h.queries.GetCluster(c.Request.Context(), req.ClusterID)
	if err != nil {
		apierror.Write(c, apierror.InvalidCredentials)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(cluster.PasswordHash), []byte(req.Password)); err != nil {
		apierror.Write(c, apierror.InvalidCredentials)
		return
	}

	resp, err := h.issueTokenResponse(req.ClusterID)
	if err != nil {
		log.Printf("failed to issue token: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req authRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Write(c, apierror.InvalidRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	_, err = h.queries.CreateCluster(c.Request.Context(), repo.CreateClusterParams{
		ID:           req.ClusterID,
		PasswordHash: string(hashed),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			apierror.Write(c, apierror.ClusterAlreadyExists)
			return
		}
		log.Printf("failed to create cluster: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	resp, err := h.issueTokenResponse(req.ClusterID)
	if err != nil {
		log.Printf("failed to issue token: %v", err)
		apierror.Write(c, apierror.Internal)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) issueTokenResponse(clusterID string) (authResponse, error) {
	now := time.Now()
	claims := middleware.AuthClaims{
		ClusterID: clusterID,
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
		ClusterID: clusterID,
		Token:     signed,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

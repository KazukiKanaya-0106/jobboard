package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
)

const ClusterIdContextKey = "cluster_id"

type AuthMiddleware struct {
	queries   repo.Querier
	jwtSecret []byte
}
type AuthClaims struct {
	ClusterId string `json:"cluster_id"`
	jwt.RegisteredClaims
}

func NewAuthMiddleware(query repo.Querier, jwtSecret []byte) *AuthMiddleware {
	return &AuthMiddleware{
		queries:   query,
		jwtSecret: jwtSecret,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, ok := extractBearerToken(c.GetHeader("Authorization"))
		if !ok {
			abortUnauthorized(c, "missing or invalid authorization header")
			return
		}

		claims := &AuthClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return m.jwtSecret, nil
		})
		if err != nil || !token.Valid {
			abortUnauthorized(c, "invalid token")
			return
		}

		c.Set(ClusterIdContextKey, claims.ClusterId)
		c.Next()
	}
}

func extractBearerToken(authHeader string) (string, bool) {
	const prefix = "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, prefix) {
		return "", false
	}
	return strings.TrimPrefix(authHeader, prefix), true
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
}

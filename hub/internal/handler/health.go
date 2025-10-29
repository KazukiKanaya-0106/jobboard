package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kanaya/jobboard-hub/internal/database"
)

type HealthHandler struct {
	db  *database.Database
	ctx context.Context
}

func NewHealthHandler(ctx context.Context, db *database.Database) *HealthHandler {
	return &HealthHandler{
		db:  db,
		ctx: ctx,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	if err := h.db.Ping(h.ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"error":  "database unavailable",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *HealthHandler) Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Job Board API",
		"version": "1.0.0",
	})
}

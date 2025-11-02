package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	ctx context.Context
}

func NewHealthHandler(ctx context.Context) *HealthHandler {
	return &HealthHandler{
		ctx: ctx,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *HealthHandler) Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Job Board Hub API",
		"version": "1.0.0",
	})
}

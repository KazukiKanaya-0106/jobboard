package router

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kanaya/jobboard-hub/internal/database"
	"github.com/kanaya/jobboard-hub/internal/handler"
)

func New(ctx context.Context, db *database.Database) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	healthHandler := handler.NewHealthHandler(ctx, db)
	r.GET("/health", healthHandler.Check)
	r.GET("/", healthHandler.Info)

	v1 := r.Group("/api/v1")
	{
		_ = v1 // TODO: Add API routes
	}

	return r
}

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

	api := r.Group("/api")
	{
		cluster := api.Group("/clusters")
	}

	return r
}

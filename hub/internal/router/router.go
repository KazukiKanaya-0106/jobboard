package router

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanaya/jobboard-hub/internal/database"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
	"github.com/kanaya/jobboard-hub/internal/handler"
	"github.com/kanaya/jobboard-hub/internal/middleware"
)

func New(ctx context.Context, db *database.Database, jwtSecret []byte, tokenTTL time.Duration) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	queries := repo.New(db.Pool)

	healthHandler := handler.NewHealthHandler(ctx, db)
	authHandler := handler.NewAuthHandler(queries, jwtSecret, tokenTTL)
	authMiddleware := middleware.NewAuthMiddleware(queries, jwtSecret)

	clusterHandler := handler.NewClusterHandler(queries)
	nodeHandler := handler.NewNodeHandler(queries)
	jobHandler := handler.NewJobHandler(queries)

	r.GET("/health", healthHandler.Check)
	r.GET("/", healthHandler.Info)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("/")
		protected.Use(authMiddleware.RequireAuth())
		{
			// クラスタ
			protected.GET("/clusters/me", clusterHandler.Me)

			// ノード
			protected.GET("/nodes", nodeHandler.List)
			protected.POST("/nodes", nodeHandler.Create)
			protected.DELETE("/nodes/:node_id", nodeHandler.Delete)

			// ジョブ
			protected.GET("/jobs", jobHandler.List)
			protected.GET("/jobs/:job_id", jobHandler.Get)
			protected.GET("/nodes/:node_id/jobs", jobHandler.ListByNode)
		}
	}

	return r
}

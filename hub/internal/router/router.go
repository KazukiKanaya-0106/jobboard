package router

import (
	"context"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kanaya/jobboard-hub/internal/database"
	"github.com/kanaya/jobboard-hub/internal/database/repo"
	"github.com/kanaya/jobboard-hub/internal/handler"
	"github.com/kanaya/jobboard-hub/internal/middleware"
)

func New(ctx context.Context, db *database.Database, allowedOrigins string, jwtSecret []byte, tokenTTL time.Duration) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(allowedOrigins, ","),
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept",
			"Authorization", "X-Requested-With",
		},
		AllowCredentials: true,
	}))

	queries := repo.New(db.Pool)

	healthHandler := handler.NewHealthHandler(ctx)
	authHandler := handler.NewAuthHandler(queries, jwtSecret, tokenTTL)
	authMiddleware := middleware.NewAuthMiddleware(queries, jwtSecret)

	clusterHandler := handler.NewClusterHandler(queries)
	nodeHandler := handler.NewNodeHandler(queries)
	jobHandler := handler.NewJobHandler(queries)
	jobTriggerHandler := handler.NewJobTriggerHandler(queries)

	router.GET("/health", healthHandler.Check)
	router.GET("/", healthHandler.Info)

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			// 認証
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

		jobTrigger := api.Group("/job-trigger")
		{
			// ジョブトリガー
			jobTrigger.POST("/start", jobTriggerHandler.StartJob)
			jobTrigger.POST("/finish", jobTriggerHandler.FinishJob)
		}
	}

	return router
}

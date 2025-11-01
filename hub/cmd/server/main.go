package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanaya/jobboard-hub/internal/config"
	"github.com/kanaya/jobboard-hub/internal/database"
	"github.com/kanaya/jobboard-hub/internal/router"
)

func main() {
	ctx := context.Background()

	cfg := config.Load()

	db, err := database.New(ctx, &cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database")

	r := router.New(ctx, db, cfg.Server.AllowedOrigins, []byte(cfg.Auth.JWTSecret), cfg.Auth.TokenTTL)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		db.Close()
		os.Exit(0)
	}()

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

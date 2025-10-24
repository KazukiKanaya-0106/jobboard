package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var command string
	flag.StringVar(&command, "cmd", "", "Migration command: up, down, version, force")
	flag.Parse()

	if command == "" {
		log.Fatal("Please specify a command: -cmd=up, -cmd=down, -cmd=version, -cmd=force")
	}

	// Build database URL
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	databaseURL := fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		log.Printf("Current version: %d, Dirty: %v", version, dirty)

	case "force":
		version := flag.Int("version", 0, "Version to force")
		flag.Parse()
		if *version == 0 {
			log.Fatal("Please specify version: -version=1")
		}
		if err := m.Force(*version); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		log.Printf("Forced to version %d", *version)

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

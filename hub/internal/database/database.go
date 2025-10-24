package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanaya/jobboard-hub/internal/config"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.DatabaseConfig) (*Database, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *Database) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

package postgres

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"featherweight/internal/config"
)

//go:embed migrations/0001_create_job_tables.up.sql
var createJobsTableSQL string

// Connect opens a pooled connection to Postgres
func Connect(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.DatabaseMaxConns)

	connectContext, cancel := context.WithTimeout(ctx, cfg.DatabaseConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectContext, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(connectContext); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// Migrate applies the embedded schema migration
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	migrateContext, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if _, err := pool.Exec(migrateContext, createJobsTableSQL); err != nil {
		return fmt.Errorf("apply jobs table migration: %w", err)
	}

	return nil
}

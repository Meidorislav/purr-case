package db

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func InitDatabase(ctx context.Context) (*Database, error) {
	pool, err := NewPool(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get pool database: %w", err)
	}
	return &Database{
		Pool: pool,
	}, nil
}

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}

func RunMigrations() error {
	databaseURL := os.Getenv("DATABASE_URL")
	m, err := migrate.New(
		"file://internal/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("migrate instance error: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up error: %w", err)
	}
	return nil
}

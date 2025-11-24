package storage

import (
	"context"
	"os"

	"github.com/UsatovPavel/PRAssign/internal/config"
	"github.com/UsatovPavel/PRAssign/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres() (*repository.PostgresRepo, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://pr-assignment:pr-assignment@db:5432/pr-assignment?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.HTTPClientTimeoutLong)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return repository.NewPostgresRepo(pool), nil
}

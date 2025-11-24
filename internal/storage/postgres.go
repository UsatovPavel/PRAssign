package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/UsatovPavel/PRAssign/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	pingCtx, cancel := context.WithTimeout(ctx, config.HTTPClientTimeoutMedium)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func WaitForDB(
	ctx context.Context,
	databaseURL string,
	retries int,
	delay time.Duration,
) (*pgxpool.Pool, error) {
	var lastErr error
	for i := 0; i < retries; i++ {
		pool, err := NewPool(ctx, databaseURL)
		if err == nil {
			return pool, nil
		}
		lastErr = err
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil, fmt.Errorf("database not ready after %d retries: %w", retries, lastErr)
}

func ClosePool(pool *pgxpool.Pool) {
	if pool == nil {
		return
	}
	pool.Close()
}

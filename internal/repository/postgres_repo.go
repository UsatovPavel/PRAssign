package repository

import (

    "github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
    Pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
    return &PostgresRepo{Pool: pool}
}

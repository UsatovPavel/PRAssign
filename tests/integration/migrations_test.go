package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestMigrationsApplied(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://pr-assignment:pr-assignment@db_test:5432/pr-assignment-test?sslmode=disable"
	}

	var db *sql.DB
	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			t.Fatalf("open db: %v", err)
		}
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deadline := time.Now().Add(10 * time.Second)
	for {
		if time.Now().After(deadline) {
			t.Fatalf("database did not become ready in time")
		}
		if err := db.PingContext(ctx); err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	tables := []string{
		"teams",
		"users",
		"pull_requests",
		"pull_request_reviewers",
	}

	for _, tbl := range tables {
		if !tableExists(ctx, t, db, tbl) {
			t.Fatalf("expected table %q to exist, but it does not", tbl)
		}
	}

	indexes := []string{
		"idx_users_team_active",
		"idx_pr_reviewers_pr",
		"idx_pr_reviewers_user",
	}

	for _, idx := range indexes {
		if !indexExists(ctx, t, db, idx) {
			t.Fatalf("expected index %q to exist, but it does not", idx)
		}
	}
}

func tableExists(ctx context.Context, t *testing.T, db *sql.DB, name string) bool {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables
			WHERE table_schema='public' AND table_name=$1
		);
	`
	if err := db.QueryRowContext(ctx, query, name).Scan(&exists); err != nil {
		t.Fatalf("checking table %q: %v", name, err)
	}
	return exists
}

func indexExists(ctx context.Context, t *testing.T, db *sql.DB, name string) bool {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM pg_indexes
			WHERE schemaname='public' AND indexname=$1
		);
	`
	if err := db.QueryRowContext(ctx, query, name).Scan(&exists); err != nil {
		t.Fatalf("checking index %q: %v", name, err)
	}
	return exists
}

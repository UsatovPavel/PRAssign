package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FactorialResultRow struct {
	JobID  string
	ItemID int64
	Input  int
	Output *string
	Status string
	Error  *string
}

type FactorialRepo struct {
	PostgresRepo
}

func NewFactorialRepo(pool *pgxpool.Pool) *FactorialRepo {
	return &FactorialRepo{PostgresRepo{Pool: pool}}
}

func (r *FactorialRepo) EnsureJob(ctx context.Context, jobID string, total int) error {
	const q = `
INSERT INTO factorial_jobs (job_id, total_items, created_at)
VALUES ($1, $2, now())
ON CONFLICT (job_id) DO NOTHING;
`
	_, err := r.Pool.Exec(ctx, q, jobID, total)
	return err
}

func (r *FactorialRepo) UpsertResult(ctx context.Context, row FactorialResultRow) error {
const q = `
INSERT INTO factorial_results (job_id, item_id, input, output, status, error, updated_at)
VALUES ($1,$2,$3,$4,$5,$6, now())
ON CONFLICT (job_id, item_id) DO UPDATE
SET output = EXCLUDED.output,
    status = EXCLUDED.status,
    error = EXCLUDED.error,
    updated_at = now();
`
	_, err := r.Pool.Exec(ctx, q, row.JobID, row.ItemID, row.Input, row.Output, row.Status, row.Error)
	return err
}

func (r *FactorialRepo) ListByJob(ctx context.Context, jobID string) ([]FactorialResultRow, error) {
	const q = `
SELECT job_id, item_id, input, output, status, error
FROM factorial_results
WHERE job_id = $1
ORDER BY item_id;
`
	rows, err := r.Pool.Query(ctx, q, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []FactorialResultRow
	for rows.Next() {
		var row FactorialResultRow
		if err := rows.Scan(&row.JobID, &row.ItemID, &row.Input, &row.Output, &row.Status, &row.Error); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return out, nil
}

func (r *FactorialRepo) GetJob(ctx context.Context, jobID string) (int, error) {
	const q = `SELECT total_items FROM factorial_jobs WHERE job_id = $1`
	var total int
	err := r.Pool.QueryRow(ctx, q, jobID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}


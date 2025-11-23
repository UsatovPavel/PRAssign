package repository

import (
	"context"
	"time"

	"github.com/UsatovPavel/PRAssign/internal/models"
)

type PullRequestPostgres struct {
	db *PostgresRepo
}

func NewPullRequestRepository(db *PostgresRepo) *PullRequestPostgres {
	return &PullRequestPostgres{db: db}
}

func (r *PullRequestPostgres) Create(ctx context.Context, pr models.PullRequest) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx,
		`INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
         VALUES ($1,$2,$3,$4,$5) ON CONFLICT (pull_request_id) DO NOTHING`,
		pr.PullRequestID, pr.PullRequestName, pr.AuthorID, string(pr.Status), pr.CreatedAt)
	if err != nil {
		return err
	}

	for _, uid := range pr.AssignedReviewers {
		_, err = tx.Exec(ctx,
			`INSERT INTO pull_request_reviewers (pull_request_id, user_id, assigned_at)
             VALUES ($1,$2,$3) ON CONFLICT (pull_request_id, user_id) DO NOTHING`,
			pr.PullRequestID, uid, time.Now().UTC())
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PullRequestPostgres) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	row := r.db.Pool.QueryRow(ctx,
		`SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
         FROM pull_requests WHERE pull_request_id = $1`,
		id)

	var pr models.PullRequest
	var status string

	if err := row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &status, &pr.CreatedAt, &pr.MergedAt); err != nil {
		return nil, models.NewAppError(models.NotFound, "pr not found")
	}

	if status == string(models.PRStatusMerged) {
		pr.Status = models.PRStatusMerged
	} else {
		pr.Status = models.PRStatusOpen
	}

	rows, err := r.db.Pool.Query(ctx,
		`SELECT user_id FROM pull_request_reviewers
         WHERE pull_request_id = $1 ORDER BY assigned_at`,
		id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, uid)
	}

	return &pr, nil
}

func (r *PullRequestPostgres) Update(ctx context.Context, pr models.PullRequest) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx,
		`UPDATE pull_requests SET pull_request_name=$1, author_id=$2,
         status=$3, created_at=$4, merged_at=$5 WHERE pull_request_id=$6`,
		pr.PullRequestName, pr.AuthorID, string(pr.Status), pr.CreatedAt, pr.MergedAt, pr.PullRequestID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM pull_request_reviewers WHERE pull_request_id=$1`,
		pr.PullRequestID)
	if err != nil {
		return err
	}

	for _, uid := range pr.AssignedReviewers {
		_, err = tx.Exec(ctx,
			`INSERT INTO pull_request_reviewers (pull_request_id, user_id, assigned_at)
             VALUES ($1,$2,$3)`,
			pr.PullRequestID, uid, time.Now().UTC())
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PullRequestPostgres) ListByReviewer(ctx context.Context, userID string) ([]models.PullRequest, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT pull_request_id
         FROM pull_request_reviewers
         WHERE user_id = $1
         ORDER BY assigned_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.PullRequest

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		pr, err := r.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		out = append(out, *pr)
	}

	return out, nil
}

func (r *PullRequestPostgres) ListAll(ctx context.Context) ([]models.PullRequest, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT pull_request_id FROM pull_requests ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.PullRequest

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		pr, err := r.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}

		out = append(out, *pr)
	}

	return out, nil
}

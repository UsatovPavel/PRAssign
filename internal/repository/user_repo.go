package repository

import (
	"context"

	"github.com/UsatovPavel/PRAssign/internal/models"
)

type UserPostgres struct {
	db *PostgresRepo
}

func NewUserRepository(db *PostgresRepo) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Upsert(ctx context.Context, u models.User) error {
	query := `
INSERT INTO users (user_id, username, team_name, is_active)
VALUES ($1,$2,$3,$4)
ON CONFLICT (user_id) DO UPDATE
SET username = EXCLUDED.username,
    team_name = EXCLUDED.team_name,
    is_active = EXCLUDED.is_active
`
	_, err := r.db.Pool.Exec(ctx, query, u.UserID, u.Username, u.TeamName, u.IsActive)
	return err
}

func (r *UserPostgres) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1`
	row := r.db.Pool.QueryRow(ctx, query, id)

	var u models.User
	if err := row.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
		return nil, models.NewAppError(models.NotFound, "user not found")
	}

	return &u, nil
}

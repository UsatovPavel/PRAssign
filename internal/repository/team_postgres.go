package repository

import (
	"context"

	"github.com/UsatovPavel/PRAssign/internal/models"
)

type TeamPostgres struct {
	db *PostgresRepo
}

func NewTeamRepository(db *PostgresRepo) *TeamPostgres {
	return &TeamPostgres{db: db}
}

func (r *TeamPostgres) CreateOrUpdate(ctx context.Context, t models.Team) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT DO NOTHING`,
		t.TeamName,
	)
	if err != nil {
		return err
	}
	//DELETE FROM users u WHERE u.team_name = $1 AND NOT EXISTS (SELECT 1 FROM pull_request_reviewers r WHERE r.user_id = u.user_id) слишком долго поэтому все пользователи остаются

	for _, m := range t.Members {
		_, err = tx.Exec(ctx,
			`INSERT INTO users (user_id, username, team_name, is_active)
             VALUES ($1,$2,$3,$4)
             ON CONFLICT (user_id) DO UPDATE
             SET username=EXCLUDED.username,
                 team_name=EXCLUDED.team_name,
                 is_active=EXCLUDED.is_active`,
			m.UserID, m.Username, t.TeamName, m.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *TeamPostgres) GetByName(ctx context.Context, name string) (*models.Team, error) {
	rows, err := r.db.Pool.Query(ctx,
		`SELECT user_id, username, is_active FROM users WHERE team_name = $1`,
		name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember

	for rows.Next() {
		var m models.TeamMember
		if err := rows.Scan(&m.UserID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	if len(members) == 0 {
		return nil, models.NewAppError(models.NotFound, "team not found")
	}

	return &models.Team{TeamName: name, Members: members}, nil
}

package service

import (
	"context"
	"log/slog"

	"github.com/UsatovPavel/PRAssign/internal/middleware"
	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/repository"
)

type TeamService struct {
	repo repository.TeamRepository
	l    *slog.Logger
}

func NewTeamService(repo repository.TeamRepository, l *slog.Logger) *TeamService {
	return &TeamService{repo: repo, l: l}
}

func (s *TeamService) CreateOrUpdateTeam(ctx context.Context, team models.Team) error {
	actingUser := ""
	if v := ctx.Value("user_id"); v != nil {
		if uid, ok := v.(string); ok && uid != "" {
			actingUser = uid
		}
	}
	if actingUser == "" {
		if v := ctx.Value(middleware.ContextUserID); v != nil {
			if uid, ok := v.(string); ok && uid != "" {
				actingUser = uid
			}
		}
	}
	if actingUser == "" {
		s.l.Error("team.createOrUpdate: no acting user in context")
		return models.NewAppError(models.NotFound, "forbidden")
	}

	isAdmin := false
	if v := ctx.Value("is_admin"); v != nil {
		if b, ok := v.(bool); ok {
			isAdmin = b
		}
	}
	if !isAdmin {
		if v := ctx.Value(middleware.ContextIsAdmin); v != nil {
			if b, ok := v.(bool); ok {
				isAdmin = b
			}
		}
	}

	if !isAdmin {
		memberOk := false
		for _, m := range team.Members {
			if m.UserID == actingUser {
				memberOk = true
				break
			}
		}
		if !memberOk {
			s.l.Warn("team.createOrUpdate: forbidden", "acting", actingUser, "team", team.TeamName)
			return models.NewAppError(models.NotFound, "forbidden")
		}
	}

	if err := s.repo.CreateOrUpdate(ctx, team); err != nil {
		s.l.Error("team createOrUpdate failed", "err", err, "team", team.TeamName)
		return err
	}
	return nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	t, err := s.repo.GetByName(ctx, name)
	if err != nil {
		s.l.Error("team get failed", "err", err, "team", name)
		return nil, err
	}

	actingUser := ""
	if v := ctx.Value("user_id"); v != nil {
		if uid, ok := v.(string); ok && uid != "" {
			actingUser = uid
		}
	}
	if actingUser == "" {
		if v := ctx.Value(middleware.ContextUserID); v != nil {
			if uid, ok := v.(string); ok && uid != "" {
				actingUser = uid
			}
		}
	}
	if actingUser == "" {
		s.l.Error("team.get: no acting user in context")
		return nil, models.NewAppError(models.NotFound, "forbidden")
	}

	isAdmin := false
	if v := ctx.Value("is_admin"); v != nil {
		if b, ok := v.(bool); ok {
			isAdmin = b
		}
	}
	if !isAdmin {
		if v := ctx.Value(middleware.ContextIsAdmin); v != nil {
			if b, ok := v.(bool); ok {
				isAdmin = b
			}
		}
	}

	if !isAdmin {
		memberOk := false
		for _, m := range t.Members {
			if m.UserID == actingUser {
				memberOk = true
				break
			}
		}
		if !memberOk {
			s.l.Warn("team.get: forbidden", "acting", actingUser, "team", name)
			return nil, models.NewAppError(models.NotFound, "forbidden")
		}
	}

	return t, nil
}

package service

import (
	"context"
	"log/slog"

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
	return t, nil
}

package service

import (
	"context"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/repository"
)

type TeamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) CreateOrUpdateTeam(ctx context.Context, team models.Team) error {
	return s.repo.CreateOrUpdate(ctx, team)
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	return s.repo.GetByName(ctx, name)
}

package service

import (
	"context"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/repository"
)

type UserService struct {
	users repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{users: repo}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	u.IsActive = isActive

	err = s.users.Upsert(ctx, *u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.users.GetByID(ctx, id)
}

func (s *UserService) GetReview(ctx context.Context, userID string) (*models.ReviewResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.ReviewResponse{
		UserID: user.UserID,
		Count:  0,
		PRs:    []models.PullRequest{},
	}, nil
}

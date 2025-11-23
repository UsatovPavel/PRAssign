package service

import (
	"context"
	"log/slog"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/repository"
)

type UserService struct {
	users repository.UserRepository
	l     *slog.Logger
}

func NewUserService(repo repository.UserRepository, l *slog.Logger) *UserService {
	return &UserService{users: repo, l: l}
}

func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	u, err := s.users.GetByID(ctx, userID)
	if err != nil {
		s.l.Error("setIsActive: get user failed", "err", err, "user", userID)
		return nil, err
	}

	u.IsActive = isActive

	if err := s.users.Upsert(ctx, *u); err != nil {
		s.l.Error("setIsActive: upsert failed", "err", err, "user", userID, "is_active", isActive)
		return nil, err
	}

	return u, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	u, err := s.users.GetByID(ctx, id)
	if err != nil {
		s.l.Error("get user by id failed", "err", err, "user", id)
		return nil, err
	}
	return u, nil
}

func (s *UserService) GetReview(ctx context.Context, userID string) (*models.ReviewResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		s.l.Error("get review: get user failed", "err", err, "user", userID)
		return nil, err
	}

	resp := &models.ReviewResponse{
		UserID: user.UserID,
		Count:  0,
		PRs:    []models.PullRequest{},
	}

	return resp, nil
}

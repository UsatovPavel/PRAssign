package service

import (
	"context"
	"log/slog"

	"github.com/UsatovPavel/PRAssign/internal/middleware"
	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/repository"
)

type UserService struct {
	users  repository.UserRepository
	prRepo repository.PullRequestRepository
	l      *slog.Logger
}

func NewUserService(userRepo repository.UserRepository, prRepo repository.PullRequestRepository, l *slog.Logger) *UserService {
	return &UserService{users: userRepo, prRepo: prRepo, l: l}
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
	actingUserInf := ctx.Value(middleware.ContextUserID)
	if actingUserInf == nil {
		s.l.Error("get review: no acting user in context")
		return nil, models.NewAppError(models.NotFound, "forbidden")
	}
	actingUser, ok := actingUserInf.(string)
	if !ok {
		s.l.Error("get review: acting user value has wrong type")
		return nil, models.NewAppError(models.NotFound, "forbidden")
	}

	isAdmin := false
	if a := ctx.Value(middleware.ContextIsAdmin); a != nil {
		if b, ok := a.(bool); ok {
			isAdmin = b
		}
	}

	if actingUser != userID && !isAdmin {
		s.l.Warn("get review: forbidden", "acting", actingUser, "target", userID)
		return nil, models.NewAppError(models.NotFound, "forbidden")
	}

	prs, err := s.prRepo.ListByReviewer(ctx, userID)
	if err != nil {
		s.l.Error("get review: repo failed", "err", err, "user", userID)
		return nil, err
	}

	resp := &models.ReviewResponse{
		UserID: userID,
		Count:  len(prs),
		PRs:    prs,
	}

	return resp, nil
}

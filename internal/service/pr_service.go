package service

import (
	"context"
	"time"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/repository"
)

type PRService struct {
	prRepo   repository.PullRequestRepository
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewPRService(
	pr repository.PullRequestRepository,
	team repository.TeamRepository,
	user repository.UserRepository,
) *PRService {
	return &PRService{
		prRepo:   pr,
		teamRepo: team,
		userRepo: user,
	}
}

func (s *PRService) Create(
	ctx context.Context,
	id string,
	name string,
	authorID string,
) (*models.PullRequest, error) {
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	team, err := s.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	var reviewers []string
	for _, m := range team.Members {
		if m.UserID != authorID && m.IsActive {
			reviewers = append(reviewers, m.UserID)
		}
	}

	now := time.Now().UTC()
	pr := models.PullRequest{
		PullRequestID:     id,
		PullRequestName:   name,
		AuthorID:          authorID,
		Status:            models.PRStatusOpen,
		AssignedReviewers: reviewers,
		CreatedAt:         &now,
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

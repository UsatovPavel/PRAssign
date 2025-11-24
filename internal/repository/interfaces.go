package repository

import (
	"context"

	"github.com/UsatovPavel/PRAssign/internal/models"
)

type UserRepository interface {
	Upsert(ctx context.Context, u models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
}

type TeamRepository interface {
	CreateOrUpdate(ctx context.Context, t models.Team) error
	GetByName(ctx context.Context, name string) (*models.Team, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pr models.PullRequest) error
	GetByID(ctx context.Context, id string) (*models.PullRequest, error)
	Update(ctx context.Context, pr models.PullRequest) error
	ListByReviewer(ctx context.Context, userID string) ([]models.PullRequest, error)
	ListAll(ctx context.Context) ([]models.PullRequest, error)
}

package service

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/repository"
)

type PRService struct {
	prRepo   repository.PullRequestRepository
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
	rnd      *rand.Rand
	l        *slog.Logger
}

func NewPullRequestService(
	pr repository.PullRequestRepository,
	team repository.TeamRepository,
	user repository.UserRepository,
	l *slog.Logger,
) *PRService {
	return &PRService{
		prRepo:   pr,
		teamRepo: team,
		userRepo: user,
		rnd:      rand.New(rand.NewSource(time.Now().UnixNano())),
		l:        l,
	}
}

func (s *PRService) Create(ctx context.Context, id, name, authorID string) (*models.PullRequest, error) {
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		s.l.Error("create pr: author not found", "err", err, "author_id", authorID, "pr_id", id)
		return nil, models.NewAppError(models.NotFound, "author not found")
	}

	team, err := s.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		s.l.Error("create pr: team not found", "err", err, "team", author.TeamName, "pr_id", id)
		return nil, models.NewAppError(models.NotFound, "team not found")
	}

	candidates := make([]string, 0, len(team.Members))
	for _, m := range team.Members {
		if m.UserID == authorID {
			continue
		}
		if !m.IsActive {
			continue
		}
		candidates = append(candidates, m.UserID)
	}

	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

	assigned := make([]string, 0, 2)
	for i := 0; i < len(candidates) && i < 2; i++ {
		assigned = append(assigned, candidates[i])
	}

	now := time.Now().UTC()
	pr := models.PullRequest{
		PullRequestID:     id,
		PullRequestName:   name,
		AuthorID:          authorID,
		Status:            models.PRStatusOpen,
		AssignedReviewers: assigned,
		CreatedAt:         &now,
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		s.l.Error("create pr: repo create failed", "err", err, "pr_id", id, "assigned", assigned)
		return nil, err
	}

	return &pr, nil
}

func (s *PRService) Merge(ctx context.Context, id string) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		s.l.Error("merge pr: get by id failed", "err", err, "pr_id", id)
		return nil, models.NewAppError(models.NotFound, "pr not found")
	}

	if pr.Status == models.PRStatusMerged {
		return pr, nil
	}

	now := time.Now().UTC()
	pr.Status = models.PRStatusMerged
	pr.MergedAt = &now

	if err := s.prRepo.Update(ctx, *pr); err != nil {
		s.l.Error("merge pr: update failed", "err", err, "pr_id", id)
		return nil, err
	}

	return pr, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldUserID string) (string, *models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		s.l.Error("reassign: get pr failed", "err", err, "pr_id", prID)
		return "", nil, models.NewAppError(models.NotFound, "pr not found")
	}

	if pr.Status == models.PRStatusMerged {
		s.l.Error("reassign: pr merged", "pr_id", prID)
		return "", nil, models.NewAppError(models.PRMerged, "cannot reassign on merged PR")
	}

	assignedIndex := -1
	for i, uid := range pr.AssignedReviewers {
		if uid == oldUserID {
			assignedIndex = i
			break
		}
	}
	if assignedIndex == -1 {
		s.l.Error("reassign: old user not assigned", "pr_id", prID, "old_user", oldUserID)
		return "", nil, models.NewAppError(models.NotAssigned, "reviewer is not assigned to this PR")
	}

	oldUser, err := s.userRepo.GetByID(ctx, oldUserID)
	if err != nil {
		s.l.Error("reassign: get old user failed", "err", err, "old_user", oldUserID)
		return "", nil, models.NewAppError(models.NotFound, "old user not found")
	}

	team, err := s.teamRepo.GetByName(ctx, oldUser.TeamName)
	if err != nil {
		s.l.Error("reassign: get team failed", "err", err, "team", oldUser.TeamName)
		return "", nil, models.NewAppError(models.NotFound, "team not found")
	}

	exclude := map[string]struct{}{}
	exclude[oldUserID] = struct{}{}
	exclude[pr.AuthorID] = struct{}{}
	for _, u := range pr.AssignedReviewers {
		exclude[u] = struct{}{}
	}

	candidates := make([]string, 0, len(team.Members))
	for _, m := range team.Members {
		if _, ok := exclude[m.UserID]; ok {
			continue
		}
		if !m.IsActive {
			continue
		}
		candidates = append(candidates, m.UserID)
	}

	if len(candidates) == 0 {
		s.l.Error("reassign: no candidate", "pr_id", prID)
		return "", nil, models.NewAppError(models.NoCandidate, "no active replacement candidate in team")
	}

	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	newUser := candidates[0]

	pr.AssignedReviewers[assignedIndex] = newUser

	if err := s.prRepo.Update(ctx, *pr); err != nil {
		s.l.Error("reassign: update pr failed", "err", err, "pr_id", prID, "new_user", newUser)
		return "", nil, err
	}

	return newUser, pr, nil
}

func (s *PRService) ListAll(ctx context.Context) ([]models.PullRequest, error) {
	out, err := s.prRepo.ListAll(ctx)
	if err != nil {
		s.l.Error("list all pr: repo failed", "err", err)
		return nil, err
	}
	return out, nil
}

func (s *PRService) ListByReviewer(ctx context.Context, userID string) ([]models.PullRequest, error) {
	out, err := s.prRepo.ListByReviewer(ctx, userID)
	if err != nil {
		s.l.Error("list by reviewer: repo failed", "err", err, "user", userID)
		return nil, err
	}
	return out, nil
}

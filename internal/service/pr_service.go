package service

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"github.com/UsatovPavel/PRAssign/internal/config"
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

func (s *PRService) Create(
	ctx context.Context,
	id, name, authorID string,
) (*models.PullRequest, error) {
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

	rand.Shuffle(
		len(candidates),
		func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] },
	)

	assigned := make([]string, 0, config.PRDefaultReviewersCap)
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

func (s *PRService) Merge(
	ctx context.Context,
	id, actingUser string,
	isAdmin bool,
) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		s.l.Error("merge pr: get by id failed", "err", err, "pr_id", id)
		return nil, models.NewAppError(models.NotFound, "pr not found")
	}

	if pr.Status == models.PRStatusMerged {
		return pr, nil
	}

	if pr.AuthorID != actingUser && !isAdmin {
		s.l.Warn("merge pr: forbidden", "acting", actingUser, "pr_author", pr.AuthorID, "pr_id", id)
		return nil, models.NewAppError(models.NotFound, "not allowed to merge")
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

func (s *PRService) ReassignReviewer(
	ctx context.Context,
	prID, oldUserID, actingUser string,
	isAdmin bool,
) (string, *models.PullRequest, error) {
	pr, err := s.loadAndCheckPR(ctx, prID, actingUser, isAdmin)
	if err != nil {
		return "", nil, err
	}

	idx, err := s.findAssignedIndex(pr, oldUserID, prID)
	if err != nil {
		return "", nil, err
	}

	newUser, err := s.pickReplacement(ctx, pr, oldUserID)
	if err != nil {
		return "", nil, err
	}

	pr.AssignedReviewers[idx] = newUser
	if err := s.prRepo.Update(ctx, *pr); err != nil {
		s.l.Error("reassign: update pr failed", "err", err, "pr_id", prID, "new", newUser)
		return "", nil, err
	}

	return newUser, pr, nil
}

func (s *PRService) loadAndCheckPR(
	ctx context.Context,
	prID, actingUser string,
	isAdmin bool,
) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, models.NewAppError(models.NotFound, "pr not found")
	}
	if pr.AuthorID != actingUser && !isAdmin {
		return nil, models.NewAppError(models.NotFound, "not allowed to reassign")
	}
	if pr.Status == models.PRStatusMerged {
		return nil, models.NewAppError(models.PRMerged, "cannot reassign on merged PR")
	}
	return pr, nil
}

func (s *PRService) findAssignedIndex(
	pr *models.PullRequest,
	oldUserID, _ string,
) (int, error) {
	for i, uid := range pr.AssignedReviewers {
		if uid == oldUserID {
			return i, nil
		}
	}
	return -1, models.NewAppError(models.NotAssigned, "reviewer is not assigned to this PR")
}

func (s *PRService) pickReplacement(
	ctx context.Context,
	pr *models.PullRequest,
	oldUserID string,
) (string, error) {
	oldUser, err := s.userRepo.GetByID(ctx, oldUserID)
	if err != nil {
		return "", models.NewAppError(models.NotFound, "old user not found")
	}

	team, err := s.teamRepo.GetByName(ctx, oldUser.TeamName)
	if err != nil {
		return "", models.NewAppError(models.NotFound, "team not found")
	}

	exclude := map[string]struct{}{
		oldUserID:   {},
		pr.AuthorID: {},
	}
	for _, u := range pr.AssignedReviewers {
		exclude[u] = struct{}{}
	}

	candidates := make([]string, 0)
	for _, m := range team.Members {
		if m.IsActive {
			if _, blocked := exclude[m.UserID]; !blocked {
				candidates = append(candidates, m.UserID)
			}
		}
	}

	if len(candidates) == 0 {
		return "", models.NewAppError(models.NoCandidate, "no active replacement candidate in team")
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	return candidates[0], nil
}

func (s *PRService) ListAll(ctx context.Context) ([]models.PullRequest, error) {
	out, err := s.prRepo.ListAll(ctx)
	if err != nil {
		s.l.Error("list all pr: repo failed", "err", err)
		return nil, err
	}
	return out, nil
}

func (s *PRService) ListByReviewer(
	ctx context.Context,
	userID string,
) ([]models.PullRequest, error) {
	out, err := s.prRepo.ListByReviewer(ctx, userID)
	if err != nil {
		s.l.Error("list by reviewer: repo failed", "err", err, "user", userID)
		return nil, err
	}
	return out, nil
}

func (s *PRService) GetByID(ctx context.Context, id string) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		s.l.Error("pr get by id failed", "err", err, "pr_id", id)
		return nil, models.NewAppError(models.NotFound, "pr not found")
	}
	return pr, nil
}

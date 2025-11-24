package models

type ReviewResponse struct {
	UserID string        `json:"user_id"`
	Count  int           `json:"count"`
	PRs    []PullRequest `json:"prs"`
}

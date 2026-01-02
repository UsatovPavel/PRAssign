package pullrequest

type CreateRequest struct {
	PullRequestID   string `json:"pull_request_id"   binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id"         binding:"required"`
}

type MergeRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldUserID     string `json:"old_user_id"     binding:"required"`
}

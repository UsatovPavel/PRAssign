package models

type ErrorCode string

const (
	TeamExists  ErrorCode = "TEAM_EXISTS"
	PRExists    ErrorCode = "PR_EXISTS"
	PRMerged    ErrorCode = "PR_MERGED"
	NotAssigned ErrorCode = "NOT_ASSIGNED"
	NoCandidate ErrorCode = "NO_CANDIDATE"
	NotFound    ErrorCode = "NOT_FOUND"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *AppError) Error() string {
	return string(e.Code) + ": " + e.Message
}

func NewAppError(code ErrorCode, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}
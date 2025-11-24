package users

type SetIsActiveRequest struct {
	UserID   string `json:"user_id"   binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

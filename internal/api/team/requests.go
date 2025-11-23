package team

type TeamMember struct {
	UserID string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	IsActive bool `json:"is_active"`
}

type AddTeamRequest struct {
	TeamName string `json:"team_name" binding:"required"`
	Members []TeamMember `json:"members"`
}
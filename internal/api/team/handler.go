package team

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/response"
	"github.com/UsatovPavel/PRAssign/internal/service"
)

type Handler struct {
	Service *service.TeamService
}

func NewHandler(s *service.TeamService) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) Add(c *gin.Context) {
	var req AddTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.Service.GetTeam(c, req.TeamName); err == nil {
		appErr := models.NewAppError(models.TeamExists, "team_name already exists")
		response.WriteAppError(c, appErr)
		return
	}

	members := make([]models.TeamMember, 0, len(req.Members))
	for _, m := range req.Members {
		members = append(members, models.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	team := models.Team{
		TeamName: req.TeamName,
		Members:  members,
	}

	err := h.Service.CreateOrUpdateTeam(c, team)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"team": team})
}

func (h *Handler) Get(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_name required"})
		return
	}

	t, err := h.Service.GetTeam(c, teamName)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, t)
}

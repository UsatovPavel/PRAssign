package team

import (
	"log/slog"
	"net/http"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/response"
	"github.com/UsatovPavel/PRAssign/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.TeamService
	l       *slog.Logger
}

func NewHandler(s *service.TeamService, l *slog.Logger) *Handler {
	return &Handler{Service: s, l: l}
}

func (h *Handler) Add(c *gin.Context) {
	var req AddTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Error("team.add: bind failed", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.Service.GetTeam(c, req.TeamName); err == nil {
		appErr := models.NewAppError(models.TeamExists, "team_name already exists")
		h.l.Warn("team.add: team exists", slog.String("team", req.TeamName))
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

	teamObj := models.Team{
		TeamName: req.TeamName,
		Members:  members,
	}

	h.l.Info("team.add: creating team", slog.String("team", teamObj.TeamName), slog.Int("members", len(teamObj.Members)))

	err := h.Service.CreateOrUpdateTeam(c, teamObj)
	if err != nil {
		h.l.Error("team.add: service failed", slog.Any("err", err), slog.String("team", teamObj.TeamName))
		response.WriteAppError(c, err)
		return
	}

	h.l.Info("team.add: created", slog.String("team", teamObj.TeamName))
	c.JSON(http.StatusCreated, gin.H{"team": teamObj})
}

func (h *Handler) Get(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		h.l.Error("team.get: missing team_name")
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_name required"})
		return
	}

	h.l.Info("team.get: request", slog.String("team", teamName))
	t, err := h.Service.GetTeam(c, teamName)
	if err != nil {
		h.l.Error("team.get: service failed", slog.Any("err", err), slog.String("team", teamName))
		response.WriteAppError(c, err)
		return
	}

	h.l.Info("team.get: success", slog.String("team", teamName), slog.Int("members", len(t.Members)))
	c.JSON(http.StatusOK, t)
}

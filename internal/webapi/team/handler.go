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

func getActingUser(c *gin.Context) (string, bool) {
	uidVal, _ := c.Get("user_id")
	isAdminVal, _ := c.Get("is_admin")
	uid, _ := uidVal.(string)
	isAdmin, _ := isAdminVal.(bool)
	return uid, isAdmin
}

func (h *Handler) bindAdd(c *gin.Context) (AddTeamRequest, bool) {
	var req AddTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Error("team.add: bind failed", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return req, true
	}
	return req, false
}

func (h *Handler) isMemberOrAdmin(req AddTeamRequest, actingUser string, isAdmin bool) bool {
	if isAdmin {
		return true
	}
	for _, m := range req.Members {
		if m.UserID == actingUser {
			return true
		}
	}
	return false
}

func toModelMembers(members []Member) []models.TeamMember {
	out := make([]models.TeamMember, 0, len(members))
	for _, m := range members {
		out = append(out, models.TeamMember{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}
	return out
}

func (h *Handler) Add(c *gin.Context) {
	req, handled := h.bindAdd(c)
	if handled {
		return
	}

	actingUser, isAdmin := getActingUser(c)
	if actingUser == "" {
		h.l.Error("team.add: missing acting user")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing token"}})
		return
	}

	if !h.isMemberOrAdmin(req, actingUser, isAdmin) {
		h.l.Warn("team.add: forbidden", slog.String("acting", actingUser), slog.String("team", req.TeamName))
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "not allowed to create team"}})
		return
	}

	if _, err := h.Service.GetTeam(c.Request.Context(), req.TeamName); err == nil {
		appErr := models.NewAppError(models.TeamExists, "team_name already exists")
		h.l.Warn("team.add: team exists", slog.String("team", req.TeamName))
		response.WriteAppError(c, appErr)
		return
	}

	teamObj := models.Team{
		TeamName: req.TeamName,
		Members:  toModelMembers(req.Members),
	}

	h.l.Info("team.add: creating team", slog.String("team", teamObj.TeamName), slog.Int("members", len(teamObj.Members)))

	if err := h.Service.CreateOrUpdateTeam(c.Request.Context(), teamObj); err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "team_name required"})
		return
	}

	actingUser, isAdmin := getActingUser(c)
	if actingUser == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{"code": "UNAUTHORIZED", "message": "missing token"},
		})
		return
	}

	h.l.Info("team.get", slog.String("team", teamName))
	t, err := h.Service.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}

	allowed := isAdmin
	if !allowed {
		for _, m := range t.Members {
			if m.UserID == actingUser {
				allowed = true
				break
			}
		}
	}
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{"code": "FORBIDDEN", "message": "not allowed to view team"},
		})
		return
	}

	c.JSON(http.StatusOK, t)
}

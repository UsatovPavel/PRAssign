package statistics

import (
	"log/slog"
	"net/http"

	"github.com/UsatovPavel/PRAssign/internal/response"
	"github.com/UsatovPavel/PRAssign/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.PRService
	l   *slog.Logger
}

func NewHandler(s *service.PRService, l *slog.Logger) *Handler {
	return &Handler{svc: s, l: l}
}

func getActingUser(c *gin.Context) (string, bool) {
	uidVal, _ := c.Get("user_id")
	isAdminVal, _ := c.Get("is_admin")
	uid, _ := uidVal.(string)
	isAdmin, _ := isAdminVal.(bool)
	return uid, isAdmin
}

func (h *Handler) AssignmentsByUsers(c *gin.Context) {
	actingUser, isAdmin := getActingUser(c)
	if actingUser == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing token"}})
		return
	}
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "admin only"}})
		return
	}

	prs, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		h.l.Error("stats.by_users: svc.ListAll failed", slog.Any("err", err))
		response.WriteAppError(c, err)
		return
	}

	byUser := map[string]int{}
	for _, pr := range prs {
		for _, uid := range pr.AssignedReviewers {
			byUser[uid]++
		}
	}

	c.JSON(http.StatusOK, gin.H{"assignments_by_user": byUser})
}

func (h *Handler) AssignmentsByPullrequests(c *gin.Context) {
	actingUser, isAdmin := getActingUser(c)
	if actingUser == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing token"}})
		return
	}
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "admin only"}})
		return
	}

	prs, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		h.l.Error("stats.by_prs: svc.ListAll failed", slog.Any("err", err))
		response.WriteAppError(c, err)
		return
	}

	byPR := map[string]int{}
	for _, pr := range prs {
		byPR[pr.PullRequestID] = len(pr.AssignedReviewers)
	}

	c.JSON(http.StatusOK, gin.H{"assignments_by_pr": byPR})
}

func (h *Handler) AssignmentsForUser(c *gin.Context) {
	actingUser, isAdmin := getActingUser(c)
	if actingUser == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing token"}})
		return
	}
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "admin only"}})
		return
	}

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "user id required"}})
		return
	}

	prs, err := h.svc.ListByReviewer(c.Request.Context(), userID)
	if err != nil {
		h.l.Error("stats.for_user: svc.ListByReviewer failed", slog.Any("err", err), slog.String("user", userID))
		response.WriteAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "assignments_count": len(prs)})
}

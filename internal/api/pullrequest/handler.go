package pullrequest

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

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Error("pullrequest.create: bind failed", slog.Any("err", err), slog.String("remote", c.ClientIP()))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid json",
			},
		})
		return
	}

	h.l.Info("pullrequest.create: request", slog.String("pr_id", req.PullRequestID), slog.String("author", req.AuthorID))

	pr, err := h.svc.Create(c, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		h.l.Error("pullrequest.create: service failed", slog.Any("err", err), slog.String("pr_id", req.PullRequestID))
		response.WriteAppError(c, err)
		return
	}

	h.l.Info("pullrequest.create: success", slog.String("pr_id", pr.PullRequestID))
	c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

func (h *Handler) Merge(c *gin.Context) {
	var req MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Error("pullrequest.merge: bind failed", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid json",
			},
		})
		return
	}

	h.l.Info("pullrequest.merge: request", slog.String("pr_id", req.PullRequestID))

	pr, err := h.svc.Merge(c, req.PullRequestID)
	if err != nil {
		h.l.Error("pullrequest.merge: service failed", slog.Any("err", err), slog.String("pr_id", req.PullRequestID))
		response.WriteAppError(c, err)
		return
	}

	h.l.Info("pullrequest.merge: success", slog.String("pr_id", pr.PullRequestID))
	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (h *Handler) Reassign(c *gin.Context) {
	var req ReassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Error("pullrequest.reassign: bind failed", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid json",
			},
		})
		return
	}

	h.l.Info("pullrequest.reassign: request", slog.String("pr_id", req.PullRequestID), slog.String("old_user", req.OldUserID))

	newUser, pr, err := h.svc.ReassignReviewer(c, req.PullRequestID, req.OldUserID)
	if err != nil {
		h.l.Error("pullrequest.reassign: service failed", slog.Any("err", err), slog.String("pr_id", req.PullRequestID))
		response.WriteAppError(c, err)
		return
	}

	h.l.Info("pullrequest.reassign: success", slog.String("pr_id", pr.PullRequestID), slog.String("new_user", newUser))
	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newUser,
	})
}

package pullrequest

import (
	"log/slog"
	"net/http"

	"github.com/UsatovPavel/PRAssign/internal/models"
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

func (h *Handler) bindCreate(c *gin.Context) (CreateRequest, bool) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Error(
			"pullrequest.create: bind failed",
			slog.Any("err", err),
			slog.String("remote", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid json",
			},
		})
		return req, true
	}
	return req, false
}

func (h *Handler) requireActingUser(c *gin.Context) (string, bool, bool) {
	actingUser, isAdmin := getActingUser(c)
	if actingUser == "" {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing token"}},
		)
		return "", false, true
	}
	return actingUser, isAdmin, false
}

func (h *Handler) Create(c *gin.Context) {
	req, handled := h.bindCreate(c)
	if handled {
		return
	}

	actingUser, isAdmin, handled := h.requireActingUser(c)
	if handled {
		return
	}
	if actingUser != req.AuthorID && !isAdmin {
		h.l.Warn(
			"pullrequest.create: forbidden",
			slog.String("acting", actingUser),
			slog.String("author", req.AuthorID),
		)
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "not allowed to create PR for this author",
				},
			},
		)
		return
	}

	h.l.Info(
		"pullrequest.create: request",
		slog.String("pr_id", req.PullRequestID),
		slog.String("author", req.AuthorID),
	)

	pr, err := h.svc.Create(c, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		h.l.Error(
			"pullrequest.create: service failed",
			slog.Any("err", err),
			slog.String("pr_id", req.PullRequestID),
		)
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

	actingUser, isAdmin := getActingUser(c)
	if actingUser == "" {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "missing token"}},
		)
		return
	}

	pr, _ := h.svc.GetByID(c.Request.Context(), req.PullRequestID)
	if pr == nil {
		response.WriteAppError(c, models.NewAppError(models.NotFound, "pr not found"))
		return
	}
	if pr.AuthorID != actingUser && !isAdmin {
		h.l.Warn(
			"pullrequest.merge: forbidden",
			slog.String("acting", actingUser),
			slog.String("pr_author", pr.AuthorID),
			slog.String("pr_id", req.PullRequestID),
		)
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "not allowed to merge"}},
		)
		return
	}

	h.l.Info("pullrequest.merge: request", slog.String("pr_id", req.PullRequestID))

	pr, err := h.svc.Merge(c.Request.Context(), req.PullRequestID, actingUser, isAdmin)
	if err != nil {
		h.l.Error(
			"pullrequest.merge: service failed",
			slog.Any("err", err),
			slog.String("pr_id", req.PullRequestID),
		)
		response.WriteAppError(c, err)
		return
	}

	h.l.Info("pullrequest.merge: success", slog.String("pr_id", pr.PullRequestID))
	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (h *Handler) bindReassign(c *gin.Context) (ReassignRequest, bool) {
	var req ReassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Error("pullrequest.reassign: bind failed", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid json",
			},
		})
		return req, true
	}
	return req, false
}

func (h *Handler) requirePRAndAuthorize(c *gin.Context, prID, actingUser string, isAdmin bool) (*models.PullRequest, bool) {
	pr, err := h.svc.GetByID(c, prID)
	if err != nil {
		response.WriteAppError(c, err)
		return nil, true
	}
	if pr == nil {
		response.WriteAppError(c, models.NewAppError(models.NotFound, "pr not found"))
		return nil, true
	}
	if pr.AuthorID != actingUser && !isAdmin {
		h.l.Warn(
			"pullrequest.reassign: forbidden",
			slog.String("acting", actingUser),
			slog.String("pr_author", pr.AuthorID),
			slog.String("pr_id", prID),
		)
		c.JSON(
			http.StatusForbidden,
			gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "not allowed to reassign"}},
		)
		return nil, true
	}
	return pr, false
}

func (h *Handler) doReassign(c *gin.Context, req ReassignRequest, actingUser string, isAdmin bool) (string, *models.PullRequest, bool) {
	newUser, pr, err := h.svc.ReassignReviewer(
		c.Request.Context(),
		req.PullRequestID,
		req.OldUserID,
		actingUser,
		isAdmin,
	)
	if err != nil {
		h.l.Error(
			"pullrequest.reassign: service failed",
			slog.Any("err", err),
			slog.String("pr_id", req.PullRequestID),
		)
		response.WriteAppError(c, err)
		return "", nil, true
	}
	return newUser, pr, false
}

func (h *Handler) Reassign(c *gin.Context) {
	req, handled := h.bindReassign(c)
	if handled {
		return
	}

	actingUser, isAdmin, handled := h.requireActingUser(c)
	if handled {
		return
	}

	_, handled = h.requirePRAndAuthorize(c, req.PullRequestID, actingUser, isAdmin)
	if handled {
		return
	}

	h.l.Info(
		"pullrequest.reassign: request",
		slog.String("pr_id", req.PullRequestID),
		slog.String("old_user", req.OldUserID),
	)

	newUser, pr, handled := h.doReassign(c, req, actingUser, isAdmin)
	if handled {
		return
	}

	h.l.Info(
		"pullrequest.reassign: success",
		slog.String("pr_id", pr.PullRequestID),
		slog.String("new_user", newUser),
	)
	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newUser,
	})
}

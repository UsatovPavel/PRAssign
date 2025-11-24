package users

import (
	"log/slog"
	"net/http"

	"github.com/UsatovPavel/PRAssign/internal/response"
	"github.com/UsatovPavel/PRAssign/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.UserService
	l       *slog.Logger
}

func NewHandler(s *service.UserService, l *slog.Logger) *Handler {
	return &Handler{service: s, l: l}
}

func (h *Handler) SetIsActive(c *gin.Context) {
	var req SetIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.l.Error("users.setIsActive: bind failed", slog.Any("err", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.l.Info("users.setIsActive: request", slog.String("user", req.UserID), slog.Bool("is_active", req.IsActive))

	u, err := h.service.SetIsActive(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		h.l.Error("users.setIsActive: service failed", slog.Any("err", err), slog.String("user", req.UserID))
		response.WriteAppError(c, err)
		return
	}

	h.l.Info("users.setIsActive: success", slog.String("user", u.UserID), slog.Bool("is_active", u.IsActive))
	c.JSON(http.StatusOK, gin.H{"user": u})
}

func (h *Handler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		h.l.Error("users.getReview: missing user_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "user_id required"}})
		return
	}
	h.l.Info("users.getReview: request", slog.String("user", userID))

	result, err := h.service.GetReview(c.Request.Context(), userID)
	if err != nil {
		h.l.Error("users.getReview: service failed", slog.Any("err", err), slog.String("user", userID))
		response.WriteAppError(c, err)
		return
	}

	h.l.Info("users.getReview: success", slog.String("user", userID), slog.Int("count", result.Count))
	c.JSON(http.StatusOK, result)
}

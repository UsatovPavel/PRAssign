package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/UsatovPavel/PRAssign/internal/service"
)

type Handler struct {
	service *service.UserService
}

func NewHandler(s *service.UserService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) SetIsActive(c *gin.Context) {
	var req SetIsActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.service.SetIsActive(c.Request.Context(), req.UserID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}

func (h *Handler) GetReview(c *gin.Context) {
	userID := c.Query("user_id")

	result, err := h.service.GetReview(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

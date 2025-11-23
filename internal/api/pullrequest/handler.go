package pullrequest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/UsatovPavel/PRAssign/internal/service"
	"github.com/UsatovPavel/PRAssign/internal/response"
)

type Handler struct {
	svc *service.PRService
}

func NewHandler(s *service.PRService) *Handler {
	return &Handler{svc: s}
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid json",
			},
		})
		return
	}

	pr, err := h.svc.Create(c, req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

func (h *Handler) Merge(c *gin.Context) {
	var req MergePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid json",
			},
		})
		return
	}

	pr, err := h.svc.Merge(c, req.PullRequestID)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (h *Handler) Reassign(c *gin.Context) {
	var req ReassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid json",
			},
		})
		return
	}

	newUser, pr, err := h.svc.ReassignReviewer(c, req.PullRequestID, req.OldUserID)
	if err != nil {
		response.WriteAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pr":          pr,
		"replaced_by": newUser,
	})
}

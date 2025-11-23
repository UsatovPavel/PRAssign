package pullrequest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/UsatovPavel/PRAssign/internal/service"
)

type Handler struct {
	Service *service.PRService
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pr, err := h.Service.CreatePR(req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"pr": pr})
}

func (h *Handler) Merge(c *gin.Context) {
	var req MergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pr, err := h.Service.MergePR(req.PullRequestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pr": pr})
}

func (h *Handler) Reassign(c *gin.Context) {
	var req ReassignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	replacedBy, pr, err := h.Service.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pr": pr, "replaced_by": replacedBy})
}
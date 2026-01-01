package factorial

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/UsatovPavel/PRAssign/internal/service"
)

type Handler struct {
	svc *service.FactorialService
}

func NewHandler(svc *service.FactorialService) *Handler {
	return &Handler{svc: svc}
}

type requestBody struct {
	Numbers []int `json:"numbers" binding:"required"`
}

type responseBody struct {
	JobID string `json:"job_id"`
	Count int    `json:"count"`
}

// POST /factorial
func (h *Handler) Enqueue(c *gin.Context) {
	var req requestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "invalid json"}})
		return
	}
	if len(req.Numbers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": "numbers must be non-empty"}})
		return
	}

	jobID := c.GetHeader("X-Job-Id")

	resp, err := h.svc.ProduceTasks(c.Request.Context(), service.FactorialRequest{
		JobID: jobID,
		Nums:  req.Numbers,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "KAFKA_ERROR", "message": err.Error()}})
		return
	}

	c.JSON(http.StatusAccepted, responseBody{
		JobID: resp.JobID,
		Count: resp.Count,
	})
}

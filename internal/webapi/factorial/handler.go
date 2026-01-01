package factorial

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/UsatovPavel/PRAssign/internal/repository"
	"github.com/UsatovPavel/PRAssign/internal/service"
)

type Handler struct {
	svc  *service.FactorialService
	repo repository.FactorialRepository
}

func NewHandler(svc *service.FactorialService, repo repository.FactorialRepository) *Handler {
	return &Handler{svc: svc, repo: repo}
}

type requestBody struct {
	Numbers []int `json:"numbers" binding:"required"`
}

type responseBody struct {
	JobID string `json:"job_id"`
	Count int    `json:"count"`
}

type resultItemResponse struct {
	ItemID int64   `json:"item_id"`
	Input  int     `json:"input"`
	Status string  `json:"status"`
	Output *string `json:"output,omitempty"`
	Error  *string `json:"error,omitempty"`
}

type resultResponse struct {
	JobID        string               `json:"job_id"`
	Status       string               `json:"status"` // done | partial
	TotalItems   int                  `json:"total_items"`
	DoneItems    int                  `json:"done_items"`
	FailedItems  int                  `json:"failed_items"`
	PendingItems int                  `json:"pending_items"`
	Items        []resultItemResponse `json:"items"`
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

	if err := h.repo.EnsureJob(c.Request.Context(), jobID, len(req.Numbers)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "DB_ERROR", "message": err.Error()}})
		return
	}

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

// GET /factorial/:job_id/result
func (h *Handler) GetResult(c *gin.Context) {
	jobID := c.Param("job_id")
	ctx := c.Request.Context()

	total, err := h.repo.GetJob(ctx, jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "job not found"}})
		return
	}

	rows, err := h.repo.ListByJob(ctx, jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "DB_ERROR", "message": err.Error()}})
		return
	}

	done := 0
	failed := 0
	items := make([]resultItemResponse, 0, len(rows))
	for _, r := range rows {
		switch r.Status {
		case "done":
			done++
		case "failed":
			failed++
		}
		item := resultItemResponse{
			ItemID: r.ItemID,
			Input:  r.Input,
			Status: r.Status,
			Output: r.Output,
			Error:  r.Error,
		}
		items = append(items, item)
	}

	processed := done + failed
	status := "partial"
	if processed == total {
		status = "done"
	}

	resp := resultResponse{
		JobID:        jobID,
		Status:       status,
		TotalItems:   total,
		DoneItems:    done,
		FailedItems:  failed,
		PendingItems: total - processed,
		Items:        items,
	}

	c.JSON(http.StatusOK, resp)
}

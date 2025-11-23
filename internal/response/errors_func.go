package response

import (
	"net/http"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/gin-gonic/gin"
)

func WriteAppError(c *gin.Context, err error) {
	if app, ok := err.(*models.AppError); ok {
		c.JSON(codeToStatus(app.Code), gin.H{
			"error": gin.H{
				"code":    app.Code,
				"message": app.Message,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    "INTERNAL",
			"message": err.Error(),
		},
	})
}

func WriteOK(c *gin.Context, payload any) {
	c.JSON(http.StatusOK, gin.H{
		"data": payload,
	})
}

func codeToStatus(code models.ErrorCode) int {
	switch code {
	case models.TeamExists,
		models.PRExists:
		return http.StatusConflict

	case models.NotFound:
		return http.StatusNotFound

	case models.NotAssigned,
		models.PRMerged,
		models.NoCandidate:
		return http.StatusConflict

	default:
		return http.StatusBadRequest
	}
}

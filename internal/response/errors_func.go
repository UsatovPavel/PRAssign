package response

import (
	"errors"
	"net/http"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/gin-gonic/gin"
)

func WriteAppError(c *gin.Context, err error) {
	var app *models.AppError
	if errors.As(err, &app) {
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

func WriteValidationError(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_REQUEST", "message": msg}})
}

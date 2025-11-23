package health

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	l *slog.Logger
}

func NewHandler(l *slog.Logger) *Handler {
	return &Handler{l: l}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	h.l.Info("health check", slog.String("remote", c.ClientIP()))
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SkipK6Logger(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ua := c.Request.UserAgent()

		if strings.Contains(strings.ToLower(ua), "k6") || c.GetHeader("X-K6-Test") == "1" {
			c.Next()
			return
		}

		reqID := uuid.New().String()
		start := time.Now()

		c.Next()

		latency := time.Since(start)

		l.With(
			slog.String("request_id", reqID),
			slog.String("ip", c.ClientIP()),
			slog.String("user_agent", ua),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.RequestURI()),
			slog.Int("status", c.Writer.Status()),
			slog.Int64("latency_ms", latency.Milliseconds()),
		).Info("request processed")
	}
}

package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SkipK6Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		ua := c.Request.UserAgent()
		if strings.Contains(strings.ToLower(ua), "k6") || c.GetHeader("X-K6-Test") == "1" {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.RequestURI()

		log.Printf("%s - \"%s\" %s %s %d %s", clientIP, ua, method, path, status, latency)
	}
}

package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		slog.Info("HTTP Request",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.String("duration", duration.String()),
			slog.String("client_ip", c.ClientIP()),
		)
	}
}

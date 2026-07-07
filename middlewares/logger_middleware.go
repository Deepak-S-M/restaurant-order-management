package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errors := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if query != "" {
			path = path + "?" + query
		}

		attrs := []slog.Attr{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
		}

		if errors != "" {
			attrs = append(attrs, slog.String("errors", errors))
		}

		switch {
		case status >= 500:
			slog.LogAttrs(c.Request.Context(), slog.LevelError, "Server Error", attrs...)
		case status >= 400:
			slog.LogAttrs(c.Request.Context(), slog.LevelWarn, "Client Error", attrs...)
		default:
			slog.LogAttrs(c.Request.Context(), slog.LevelInfo, "Request", attrs...)
		}
	}
}

package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/notFil/cspotlight/pkg/logger"
	"go.uber.org/zap"
)

// RequestLogger adds a logger to the context with a request ID
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Generate or get Request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)

		// Create a logger with the request ID
		// We use the global logger as the base
		if logger.Logger == nil {
			// Fallback if logger wasn't initialized
			logger.InitializeLogger("development")
		}

		log := logger.Logger.With(zap.String("request_id", requestID))

		// Set the logger in the context
		c.Set("logger", log)

		// Process request
		c.Next()

		// Log request details
		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Log errors if any
			for _, e := range c.Errors.Errors() {
				log.Error(e)
			}
		} else {
			log.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
			)
		}
	}
}

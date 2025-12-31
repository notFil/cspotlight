package middleware

import (
	"time"

	"cspotlight/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestLogger adds a logger to the context with a request ID
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)

		log := logger.Logger.With(zap.String("request_id", requestID))

		// Attach logger to Go context
		ctx := logger.WithLogger(c.Request.Context(), log)
		c.Request = c.Request.WithContext(ctx)

		// Continue the request
		c.Next()

		latency := time.Since(start)

		log.Info("request completed",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		)
	}
}

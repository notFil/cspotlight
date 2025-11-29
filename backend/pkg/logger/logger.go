package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Logger *zap.Logger

const loggerKey = "logger"

func InitializeLogger(env string) {
	if env == "production" {
		Logger, _ = zap.NewProduction()
	} else {
		Logger, _ = zap.NewDevelopment()
	}
}

// FromContext returns the logger associated with the context.
// If no logger is found, it returns the global logger.
func FromContext(ctx context.Context) *zap.Logger {
	if c, ok := ctx.(*gin.Context); ok {
		if l, exists := c.Get(loggerKey); exists {
			if logger, ok := l.(*zap.Logger); ok {
				return logger
			}
		}
	}

	// Fallback for standard context
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return l
	}

	return Logger
}

package logger

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Logger *zap.Logger

type loggerKey struct{}

func InitializeLogger(isProduction bool) {
	if isProduction {
		Logger, _ = zap.NewProduction()
	} else {
		Logger, _ = zap.NewDevelopment()
	}
}

// WithLogger attaches a logger to the context
func WithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// FromContext returns the logger associated with the context.
// If no logger is found, it returns the global logger.
func FromContext(ctx context.Context) *zap.Logger {
	if c, ok := ctx.(*gin.Context); ok {
		if l, exists := c.Get(loggerKey{}); exists {
			if logger, ok := l.(*zap.Logger); ok {
				return logger
			}
		}
	}

	if l, ok := ctx.Value(loggerKey{}).(*zap.Logger); ok {
		return l
	}

	return Logger
}

package middleware

import (
	"net/http"
	"sync"

	apperrors "cspotlight/internal/errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(15, 15)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.Error(apperrors.New(http.StatusTooManyRequests, "Too many requests. Please try again later."))
			c.Abort()
			return
		}
	}
}

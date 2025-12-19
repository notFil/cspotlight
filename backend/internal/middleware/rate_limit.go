package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"golang.org/x/time/rate"
)

var limiters = make(map[string]*rate.Limiter)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if _, exists := limiters[ip]; !exists {
			limiters[ip] = rate.NewLimiter(10, 10)
		}

		limiter := limiters[ip]
		if !limiter.Allow() {
			c.Error(apperrors.New(http.StatusTooManyRequests, "Too many requests. Please try again later."))
			return
		}
	}
}

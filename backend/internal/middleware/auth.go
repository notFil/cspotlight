package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/notFil/cspotlight/internal/auth"

	"github.com/notFil/cspotlight/internal/logger"

	"github.com/notFil/cspotlight/internal/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorResponse(c, http.StatusUnauthorized, "authorization header required")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := auth.ValidateAccessToken(tokenString, secretKey)
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		ctx = auth.ContextWithClaims(ctx, claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequiredRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.FromContext(c)

		ctx := c.Request.Context()

		claims := auth.GetUserClaims(ctx)
		if claims == nil {
			logger.Error("claims not found in context")
			response.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
			c.Abort()
			return
		}

		if isAuthorized := slices.Contains(allowedRoles, claims.Role); !isAuthorized {
			logger.Error("insufficient permissions")
			response.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

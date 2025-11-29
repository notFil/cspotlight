package middleware

import (
	"net/http"
	"strings"

	"github.com/notFil/cspotlight/pkg/constants"

	"github.com/notFil/cspotlight/pkg/auth"

	"github.com/notFil/cspotlight/pkg/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := auth.ValidateAccessToken(tokenString, secretKey)
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		}

		c.Set(constants.ClaimsContextKey, claims)

		c.Next()
	}
}

func RequiredRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("Role")
		if !exists {
			response.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			response.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		isAuthorized := false
		for _, role := range allowedRoles {
			if roleStr == role {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			response.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/common"
)

func JWTAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := common.ValidateJWT(secretKey, tokenString)
		if err != nil {
			common.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		}
		c.Set("UserID", claims.Subject)
		c.Set("TeamID", claims.TeamID)
		c.Set("Role", claims.Role)

		c.Next()
	}
}

func RequiredRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("Role")
		if !exists {
			common.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			common.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
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
			common.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetClaims() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("UserID")
		teamID := c.MustGet("TeamID")
		role := c.MustGet("Role")

		teamIDStr := teamID.(string)
		claims := common.AuthClaims{
			UserID: userID.(string),
			TeamID: &teamIDStr,
			Role:   role.(string),
		}

		c.Set(common.ClaimsContextKey, claims)
		c.Next()
	}
}

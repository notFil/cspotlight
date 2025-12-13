package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/notFil/cspotlight/internal/auth"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	logger "github.com/notFil/cspotlight/internal/logger"
)

func AuthMiddleware(c *gin.Context) {
	logger := logger.FromContext(c)

	session := sessions.Default(c)
	userID := session.Get("userID")
	teamID := session.Get("teamID")
	role := session.Get("role")

	if userID == nil || teamID == nil || role == nil {
		logger.Error("unauthorized")
		c.Error(apperrors.New(http.StatusUnauthorized, "unauthorized"))
		c.Abort()
		return
	}

	userContext := auth.NewUserContext(userID.(string), teamID.(string), role.(string))

	ctx := auth.ContextWithUser(c, userContext)

	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

func RequiredRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.FromContext(c)

		session := sessions.Default(c)
		role := session.Get("role")

		if isAuthorized := slices.Contains(allowedRoles, role.(string)); !isAuthorized {
			logger.Error("insufficient permissions")
			c.Error(apperrors.New(http.StatusForbidden, "insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

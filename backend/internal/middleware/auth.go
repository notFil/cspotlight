package middleware

import (
	"net/http"
	"slices"

	"cspotlight/internal/auth"
	apperrors "cspotlight/internal/errors"
	logger "cspotlight/internal/logger"
	"cspotlight/internal/util"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(c *gin.Context) {
	logger := logger.FromContext(c)

	session := sessions.Default(c)
	userID := session.Get("userID")
	teamID := session.Get("teamID")
	role := session.Get("role")

	if userID == nil || role == nil {
		logger.Error("unauthorized")
		c.Error(apperrors.New(http.StatusUnauthorized, "unauthorized"))
		c.Abort()
		return
	}

	parsedUserID, _ := util.ParseUUIDValue(userID)
	parsedTeamID, _ := util.ParseUUIDValue(teamID)

	userContext := auth.NewUserContext(parsedUserID, parsedTeamID, role.(string))

	logger.Info("user context created", zap.String("userID", userContext.ID.String()), zap.String("teamID", userContext.TeamID.String()), zap.String("role", userContext.Role))

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

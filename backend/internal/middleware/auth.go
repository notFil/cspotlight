package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/auth"
	apperrors "github.com/notFil/cspotlight/internal/errors"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/store"

	"github.com/gin-gonic/gin"
)

func JWTAuth(cache store.Cache, secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		tokenString, err := auth.ParseAuthHeader(authHeader)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		claims, err := auth.ValidateAccessToken(tokenString, secretKey)
		if err != nil {
			c.Error(apperrors.New(http.StatusUnauthorized, err.Error()))
			c.Abort()
			return
		}

		subject := uuid.MustParse(claims.Subject)
		teamID := uuid.MustParse(claims.TeamID)

		_, err = cache.Get(c.Request.Context(), fmt.Sprintf("user:%s:%s", subject, claims.ID))
		if err != nil {
			c.Error(apperrors.New(http.StatusUnauthorized, "invalid token"))
			c.Abort()
			return
		}

		user := auth.AuthContext{
			Subject: subject,
			TeamID:  teamID,
			Role:    claims.Role,
		}

		ctx := c.Request.Context()
		ctx = auth.ContextWithUser(ctx, &user)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func RequiredRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logger.FromContext(c)

		ctx := c.Request.Context()

		user := auth.GetUserContext(ctx)
		if user == nil {
			c.Error(apperrors.New(http.StatusForbidden, "insufficient permissions"))
			c.Abort()
			return
		}

		if isAuthorized := slices.Contains(allowedRoles, user.Role); !isAuthorized {
			logger.Error("insufficient permissions")
			c.Error(apperrors.New(http.StatusForbidden, "insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

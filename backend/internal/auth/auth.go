package auth

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/notFil/cspotlight/internal/constants"
)

type contextKey struct{}

type UserContext struct {
	ID     uuid.UUID
	TeamID uuid.UUID
	Role   string
}

func NewUserContext(id string, teamID string, role string) *UserContext {
	parsedID := uuid.MustParse(id)
	parsedTeamID := uuid.MustParse(teamID)

	return &UserContext{
		ID:     parsedID,
		TeamID: parsedTeamID,
		Role:   role,
	}
}

func (c *UserContext) IsAdmin() bool {
	return strings.EqualFold(c.Role, "admin")
}

func (c *UserContext) IsSuperadmin() bool {
	return strings.EqualFold(c.Role, "superadmin")
}

func (c *UserContext) IsSameUser(id uuid.UUID) bool {
	return c.ID == id
}

func (c *UserContext) HasSameTeam(id uuid.UUID) bool {
	return c.TeamID == id
}

func ContextWithUser(ctx context.Context, userContext *UserContext) context.Context {
	return context.WithValue(ctx, contextKey{}, userContext)
}

func GetUserContext(ctx context.Context) *UserContext {
	if userContext, ok := ctx.Value(contextKey{}).(*UserContext); ok {
		return userContext
	}

	if val := ctx.Value(constants.UserContextKey); val != nil {
		if userContext, ok := val.(*UserContext); ok {
			return userContext
		}
	}

	return nil
}

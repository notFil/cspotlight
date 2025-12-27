package auth

import (
	"context"
	"strings"

	"cspotlight/internal/constants"

	"github.com/google/uuid"
)

type contextKey struct{}

type UserContext struct {
	ID     uuid.UUID
	TeamID uuid.UUID
	Role   string
}

func NewUserContext(id uuid.UUID, teamID uuid.UUID, role string) *UserContext {

	return &UserContext{
		ID:     id,
		TeamID: teamID,
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

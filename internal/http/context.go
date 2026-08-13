package http

import (
	"context"

	"github.com/google/uuid"
)

// need to get User from context, user from request

type contextKey string

const userContextKey contextKey = "user"

type UserContext struct {
	UUID  uuid.UUID
	Email string
	Role  string
}

func ContextWithUser(ctx context.Context, user *UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func UserFromContext(ctx context.Context) *UserContext {
	u, ok := ctx.Value(userContextKey).(*UserContext)
	if !ok {
		return nil
	}
	return u
}

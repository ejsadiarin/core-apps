package auth

import (
	"context"
	"core-gateway/internal/sqlc"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// context keys for user info
type contextKey string

const (
	userContextKey   contextKey = "user"
	userIDContextKey contextKey = "user_id"
)

// UserContext holds user information extracted from session
type UserContext struct {
	ID    uuid.UUID
	Email string
	Role  string
}

// AuthMiddleware extracts session cookie, validates it, and loads user into context
func AuthMiddleware(queries *sqlc.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// try to get session token from cookie
			token, err := GetSessionToken(c)
			if err != nil {
				// no session cookie - continue as guest
				return next(c)
			}

			// hash the token and look up the session
			tokenHash := HashSessionToken(token)
			session, err := queries.GetSessionByTokenHash(c.Request().Context(), tokenHash)
			if err != nil {
				// invalid or expired session - clear cookie and continue as guest
				ClearSessionCookie(c)
				return next(c)
			}

			// set user context
			userCtx := &UserContext{
				ID:    session.UserID,
				Email: session.Email,
				Role:  session.Role,
			}

			// store in echo context
			c.Set(string(userContextKey), userCtx)
			c.Set(string(userIDContextKey), session.UserID)

			return next(c)
		}
	}
}

// RequireAuth middleware returns 401 if no valid session exists
func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := GetUserFromContext(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authentication required",
				})
			}
			return next(c)
		}
	}
}

// RequireRole middleware checks if user has one of the allowed roles
func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := GetUserFromContext(c)
			if user == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Authentication required",
				})
			}

			// check if user's role is in allowed roles
			for _, role := range allowedRoles {
				if user.Role == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Insufficient permissions",
			})
		}
	}
}

// GetUserFromContext retrieves the user from echo context
func GetUserFromContext(c echo.Context) *UserContext {
	user, ok := c.Get(string(userContextKey)).(*UserContext)
	if !ok {
		return nil
	}
	return user
}

// GetUserIDFromContext retrieves the user ID from echo context as pgtype.UUID
func GetUserIDFromContext(c echo.Context) pgtype.UUID {
	userID, ok := c.Get(string(userIDContextKey)).(pgtype.UUID)
	if !ok {
		return pgtype.UUID{}
	}
	return userID
}

// IsAuthenticated returns true if user is logged in
func IsAuthenticated(c echo.Context) bool {
	return GetUserFromContext(c) != nil
}

// IsAdmin returns true if user has admin role
func IsAdmin(c echo.Context) bool {
	user := GetUserFromContext(c)
	return user != nil && user.Role == RoleAdmin
}

// IsGuest returns true if user is not authenticated or has guest role
func IsGuest(c echo.Context) bool {
	user := GetUserFromContext(c)
	return user == nil || user.Role == RoleGuest
}

// ContextWithUser adds user to standard context (for passing to services)
func ContextWithUser(ctx context.Context, user *UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext retrieves user from standard context
func UserFromContext(ctx context.Context) *UserContext {
	user, ok := ctx.Value(userContextKey).(*UserContext)
	if !ok {
		return nil
	}
	return user
}

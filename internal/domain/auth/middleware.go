package auth

import (
	"context"
	"encoding/json"
	"net/http"

	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"

	"github.com/google/uuid"
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
func AuthMiddleware(queries *sqlc.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := GetSessionToken(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			tokenHash := HashSessionToken(token)
			session, err := queries.GetSessionByTokenHash(r.Context(), tokenHash)
			if err != nil {
				ClearSessionCookie(w)
				next.ServeHTTP(w, r)
				return
			}

			userCtx := &UserContext{
				ID:    session.UserID,
				Email: session.Email,
				Role:  session.Role,
			}

			ctx := context.WithValue(r.Context(), userContextKey, userCtx)
			ctx = context.WithValue(ctx, userIDContextKey, session.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth middleware returns 401 if no valid session exists
func RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r)
			if user == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Authentication required",
				})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole middleware checks if user has one of the allowed roles
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := GetUserFromContext(r)
			if user == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Authentication required",
				})
				return
			}

			for _, role := range allowedRoles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Insufficient permissions",
			})
		})
	}
}

// GetUserFromContext retrieves the user from request context
func GetUserFromContext(r *http.Request) *UserContext {
	user, ok := r.Context().Value(userContextKey).(*UserContext)
	if !ok {
		return nil
	}
	return user
}

// IsAuthenticated returns true if user is logged in
func IsAuthenticated(r *http.Request) bool {
	return GetUserFromContext(r) != nil
}

// IsAdmin returns true if user has admin role
func IsAdmin(r *http.Request) bool {
	user := GetUserFromContext(r)
	return user != nil && user.Role == RoleAdmin
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

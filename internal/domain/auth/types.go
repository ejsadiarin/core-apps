package auth

import (
	"context"

	usertypes "github.com/ejsadiarin/coregateway/internal/domain/user/v1"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRespository interface {
	GetUserByEmail(ctx context.Context, user usertypes.User) (usertypes.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (usertypes.User, error)
	CreateUser(ctx context.Context, email, passwordHash, role string) (usertypes.User, error)
	CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt pgtype.Timestamp) error
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (SessionRow, error)
	DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error
	DeleteExpiredSessions(ctx context.Context) error
}

type SessionRow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt pgtype.Timestamp
	CreatedAt pgtype.Timestamp
	Email     string
	Role      string
}

// request DTOs

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

// response DTOs

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"created_at"`
}

package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*SessionRow, error)
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



package auth

import (
	"context"
	"log/slog"
	"time"

	db "github.com/ejsadiarin/coregateway/internal/db/sqlc"
	"github.com/google/uuid"
)

// type UserRespository interface {
// 	GetUserByEmail(ctx context.Context, user usertypes.User) (usertypes.User, error)
// 	GetUserByID(ctx context.Context, id uuid.UUID) (usertypes.User, error)
// 	CreateUser(ctx context.Context, email, passwordHash, role string) (usertypes.User, error)
// 	CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt pgtype.Timestamp) error
// 	GetSessionByTokenHash(ctx context.Context, tokenHash string) (SessionRow, error)
// 	DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error
// 	DeleteExpiredSessions(ctx context.Context) error
// }
//

type Service struct {
	q      db.Queries
	logger *slog.Logger
}

func (s *Service) CreateSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	// only create session if successful login
	return nil
}

func (s *Service) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*SessionRow, error) {
	return nil, nil
}

func (s *Service) DeleteSessionByTokenHash(ctx context.Context, tokenHash string) error {
	return nil
}

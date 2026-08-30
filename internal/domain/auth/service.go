package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/ejsadiarin/coregateway/internal/crypto"
	db "github.com/ejsadiarin/coregateway/internal/db/sqlc"
	usertypes "github.com/ejsadiarin/coregateway/internal/domain/user/v1"
	"github.com/ejsadiarin/coregateway/internal/session"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service defines the interface for auth operations.
type Service interface {
	Register(ctx context.Context, email, password string) (*usertypes.User, error)
	Login(ctx context.Context, email, password string) (*usertypes.User, error)
	CreateSession(ctx context.Context, userID uuid.UUID, rememberMe bool) (string, error)
	Logout(ctx context.Context, token string) error
	GetSessionUser(ctx context.Context, token string) (*usertypes.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*usertypes.User, error)
	LoginAsDemo(ctx context.Context) (*usertypes.User, string, error)
}

// authService implements Service.
type authService struct {
	queries db.Querier
	logger  *slog.Logger
}

// NewService creates a new auth service.
func NewService(queries db.Querier, logger *slog.Logger) Service {
	return &authService{
		queries: queries,
		logger:  logger,
	}
}

// Register creates a new user account.
func (s *authService) Register(ctx context.Context, email, password string) (*usertypes.User, error) {
	_, err := s.queries.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
		Role:         RoleUser,
	})
	if err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User registered", "email", email)

	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

// Login authenticates a user.
func (s *authService) Login(ctx context.Context, email, password string) (*usertypes.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.PasswordHash.Valid {
		return nil, fmt.Errorf("invalid email or password")
	}

	valid, err := crypto.VerifyPassword(password, user.PasswordHash.String)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid email or password")
	}

	s.logger.Info("User authenticated", "email", email)

	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

// CreateSession creates a new session for a user and returns the raw token.
func (s *authService) CreateSession(ctx context.Context, userID uuid.UUID, rememberMe bool) (string, error) {
	token, tokenHash, err := crypto.GenerateSessionToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	expiry := session.SessionExpiry
	if rememberMe {
		expiry = session.RememberMeExpiry
	}
	expiresAt := time.Now().Add(expiry)

	_, err = s.queries.CreateSession(ctx, db.CreateSessionParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: pgtype.Timestamp{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return token, nil
}

// Logout deletes the session identified by the given token.
func (s *authService) Logout(ctx context.Context, token string) error {
	tokenHash := crypto.HashSessionToken(token)
	return s.queries.DeleteSessionByTokenHash(ctx, tokenHash)
}

// GetSessionUser looks up a session by token hash and returns the associated user.
func (s *authService) GetSessionUser(ctx context.Context, token string) (*usertypes.User, error) {
	tokenHash := crypto.HashSessionToken(token)
	session, err := s.queries.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired session: %w", err)
	}

	user, err := s.queries.GetUser(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

// GetUser retrieves a user by ID.
func (s *authService) GetUser(ctx context.Context, id uuid.UUID) (*usertypes.User, error) {
	user, err := s.queries.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

// LoginAsDemo creates a session for the demo user.
func (s *authService) LoginAsDemo(ctx context.Context) (*usertypes.User, string, error) {
	user, err := s.queries.GetUserByEmail(ctx, DemoUserEmail)
	if err != nil {
		return nil, "", fmt.Errorf("demo user not available: %w", err)
	}

	token, err := s.CreateSession(ctx, user.ID, false)
	if err != nil {
		return nil, "", err
	}

	s.logger.Info("Demo user logged in")

	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format(time.RFC3339),
	}, token, nil
}

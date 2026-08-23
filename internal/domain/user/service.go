package user

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/ejsadiarin/coregateway/internal/crypto"
	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"
	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	usertypes "github.com/ejsadiarin/coregateway/internal/domain/user/v1"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Service defines the interface for user operations.
type Service interface {
	ListUsers(ctx context.Context) ([]*usertypes.User, error)
	CreateUser(ctx context.Context, email, password, role string) (*usertypes.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*usertypes.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, params UpdateParams, callerIsAdmin bool) (*usertypes.User, error)
	DeleteUser(ctx context.Context, id, callerID uuid.UUID, callerRole string) error
}

type UpdateParams struct {
	Email    *string
	Password *string
	Role     *string
}

// userService implements Service.
type userService struct {
	db     sqlc.Querier
	logger *slog.Logger
}

func NewService(db sqlc.Querier, logger *slog.Logger) Service {
	return &userService{db: db, logger: logger}
}

func (s *userService) ListUsers(ctx context.Context) ([]*usertypes.User, error) {
	users, err := s.db.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*usertypes.User, len(users))
	for i, u := range users {
		res[i] = &usertypes.User{
			Id:        u.ID.String(),
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: u.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return res, nil
}

func (s *userService) CreateUser(ctx context.Context, email, password, role string) (*usertypes.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	_, err := s.db.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.db.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
		Role:         role,
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info("User created", "email", email, "role", role)
	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*usertypes.User, error) {
	user, err := s.db.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, params UpdateParams, callerIsAdmin bool) (*usertypes.User, error) {
	if params.Role != nil && !callerIsAdmin {
		return nil, errors.New("only administrators can change user roles")
	}

	_, err := s.db.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	updateParams := sqlc.UpdateUserParams{ID: id}

	if params.Email != nil {
		updateParams.Email = pgtype.Text{String: strings.ToLower(strings.TrimSpace(*params.Email)), Valid: true}
	}

	if params.Password != nil {
		hashedPassword, err := crypto.HashPassword(*params.Password)
		if err != nil {
			return nil, err
		}
		updateParams.PasswordHash = pgtype.Text{String: hashedPassword, Valid: true}
	}

	if params.Role != nil {
		updateParams.Role = pgtype.Text{String: *params.Role, Valid: true}
	}

	user, err := s.db.UpdateUser(ctx, updateParams)
	if err != nil {
		return nil, err
	}

	s.logger.Info("User updated", "user_id", id.String())
	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id, callerID uuid.UUID, callerRole string) error {
	if callerID == id {
		return errors.New("cannot delete your own account")
	}

	user, err := s.db.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	if user.Email == auth.DemoUserEmail {
		return errors.New("cannot delete demo user")
	}

	if err := s.db.DeleteUserSessions(ctx, id); err != nil {
		s.logger.Warn("Failed to delete user sessions", "error", err, "user_id", id.String())
	}

	if err := s.db.DeleteUser(ctx, id); err != nil {
		return err
	}

	s.logger.Info("User deleted", "user_id", id.String())
	return nil
}

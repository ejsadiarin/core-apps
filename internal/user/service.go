package user

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/ejsadiarin/coregateway/internal/crypto"
	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"
	"github.com/ejsadiarin/coregateway/internal/auth"
	usertypes "github.com/ejsadiarin/coregateway/internal/user/v1"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateParams struct {
	Email    *string
	Password *string
	Role     *string
}

type Service struct {
	db sqlc.Querier
}

func NewService(db sqlc.Querier) *Service {
	return &Service{db: db}
}

func (s *Service) ListUsers(ctx context.Context) ([]*usertypes.User, error) {
	slog.Debug("user.Service.ListUsers")
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

func (s *Service) CreateUser(ctx context.Context, email, password, role string) (*usertypes.User, error) {
	slog.Debug("user.Service.CreateUser", "email", email, "role", role)
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

	slog.Info("user.Service.CreateUser: user created", "email", email, "role", role)
	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*usertypes.User, error) {
	slog.Debug("user.Service.GetUser", "id", id)
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

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, params UpdateParams, callerIsAdmin bool) (*usertypes.User, error) {
	slog.Debug("user.Service.UpdateUser", "id", id)
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

	slog.Info("user.Service.UpdateUser: user updated", "id", id)
	return &usertypes.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, id, callerID uuid.UUID, callerRole string) error {
	slog.Debug("user.Service.DeleteUser", "id", id)
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
		slog.Warn("user.Service.DeleteUser: failed to delete user sessions", "error", err, "id", id)
	}

	if err := s.db.DeleteUser(ctx, id); err != nil {
		return err
	}

	slog.Info("user.Service.DeleteUser: user deleted", "id", id)
	return nil
}

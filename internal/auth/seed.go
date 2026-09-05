package auth

import (
	"context"
	"log/slog"

	"github.com/ejsadiarin/coregateway/internal/crypto"
	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	DemoUserEmail = "demo@example.com"
	RoleGuest     = "guest"
	RoleUser      = "user"
	RoleAdmin     = "admin"
)

func SeedUsers(ctx context.Context, queries *sqlc.Queries, adminEmail, adminPassword string) error {
	if err := seedAdminUser(ctx, queries, adminEmail, adminPassword); err != nil {
		return err
	}

	if err := seedDemoUser(ctx, queries); err != nil {
		return err
	}

	return nil
}

func seedAdminUser(ctx context.Context, queries *sqlc.Queries, adminEmail, adminPassword string) error {
	if adminEmail == "" || adminPassword == "" {
		slog.Warn("ADMIN_EMAIL or ADMIN_PASSWORD not set, skipping admin user seed")
		return nil
	}

	_, err := queries.GetUserByEmail(ctx, adminEmail)
	if err == nil {
		slog.Debug("auth.SeedUsers: admin user already exists", "email", adminEmail)
		return nil
	}

	hashedPassword, err := crypto.HashPassword(adminPassword)
	if err != nil {
		return err
	}

	_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        adminEmail,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
		Role:         RoleAdmin,
	})
	if err != nil {
		return err
	}

	slog.Info("auth.SeedUsers: admin user created", "email", adminEmail)
	return nil
}

func seedDemoUser(ctx context.Context, queries *sqlc.Queries) error {
	_, err := queries.GetUserByEmail(ctx, DemoUserEmail)
	if err == nil {
		slog.Debug("auth.SeedUsers: demo user already exists", "email", DemoUserEmail)
		return nil
	}

	_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        DemoUserEmail,
		PasswordHash: pgtype.Text{Valid: false},
		Role:         RoleGuest,
	})
	if err != nil {
		return err
	}

	slog.Info("auth.SeedUsers: demo user created", "email", DemoUserEmail)
	return nil
}

func GetDemoUserID(ctx context.Context, queries *sqlc.Queries) (uuid.UUID, error) {
	user, err := queries.GetUserByEmail(ctx, DemoUserEmail)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

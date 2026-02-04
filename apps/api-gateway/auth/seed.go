package auth

import (
	"context"
	"core-gateway/internal/sqlc"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

const (
	DemoUserEmail = "demo@example.com"
	RoleGuest     = "guest"
	RoleUser      = "user"
	RoleAdmin     = "admin"
)

// SeedUsers creates the initial admin and demo users if they don't exist
func SeedUsers(ctx context.Context, queries *sqlc.Queries, logger *zerolog.Logger) error {
	// seed admin user from env vars
	if err := seedAdminUser(ctx, queries, logger); err != nil {
		return err
	}

	// seed demo user for guest access
	if err := seedDemoUser(ctx, queries, logger); err != nil {
		return err
	}

	return nil
}

func seedAdminUser(ctx context.Context, queries *sqlc.Queries, logger *zerolog.Logger) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		logger.Warn().Msg("ADMIN_EMAIL or ADMIN_PASSWORD not set, skipping admin user seed")
		return nil
	}

	// check if admin already exists
	_, err := queries.GetUserByEmail(ctx, adminEmail)
	if err == nil {
		logger.Debug().Str("email", adminEmail).Msg("Admin user already exists")
		return nil
	}

	// hash the password
	hashedPassword, err := HashPassword(adminPassword)
	if err != nil {
		return err
	}

	// create admin user
	_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        adminEmail,
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
		Role:         RoleAdmin,
	})
	if err != nil {
		return err
	}

	logger.Info().Str("email", adminEmail).Msg("Admin user created")
	return nil
}

func seedDemoUser(ctx context.Context, queries *sqlc.Queries, logger *zerolog.Logger) error {
	// check if demo user already exists
	_, err := queries.GetUserByEmail(ctx, DemoUserEmail)
	if err == nil {
		logger.Debug().Str("email", DemoUserEmail).Msg("Demo user already exists")
		return nil
	}

	// create demo user with no password (can't login)
	_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        DemoUserEmail,
		PasswordHash: pgtype.Text{Valid: false}, // NULL password
		Role:         RoleGuest,
	})
	if err != nil {
		return err
	}

	logger.Info().Str("email", DemoUserEmail).Msg("Demo user created")
	return nil
}

// GetDemoUserID retrieves the demo user's ID for guest access
func GetDemoUserID(ctx context.Context, queries *sqlc.Queries) (uuid.UUID, error) {
	user, err := queries.GetUserByEmail(ctx, DemoUserEmail)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

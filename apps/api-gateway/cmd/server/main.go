// @title Core Homelab API
// @version 2.0
// @description API for Core Homelab Dashboard with service monitoring and budget tracking
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @schemes http https

package main

import (
	"net/http"
	"os"
	"time"

	"core-gateway/internal/app"

	_ "core-gateway/docs"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
)

func main() {
	// initialize zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if os.Getenv("ENV") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	logger.Info().Msg("Starting Core Homelab API")

	// build configuration from environment
	cfg := app.Config{
		DatabaseURL:         getEnvOrDefault("DATABASE_URL", "postgresql://core:core@postgres:5432/core?sslmode=disable"),
		Port:                getEnvOrDefault("PORT", "8080"),
		Env:                 getEnvOrDefault("ENV", "development"),
		FrontendURL:         os.Getenv("FRONTEND_URL"),
		HealthCheckInterval: 60 * time.Second,
		AdminEmail:          os.Getenv("ADMIN_EMAIL"),
		AdminPassword:       os.Getenv("ADMIN_PASSWORD"),
	}

	// create application
	application, err := app.New(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create application")
	}
	defer application.Close()

	// start server
	if err := application.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

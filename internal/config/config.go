package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds application configuration
type Config struct {
	DatabaseURL         string
	ServerAddr          string
	Env                 string
	HealthCheckInterval time.Duration
	SessionSecret       string
	AllowedOrigins      string
	AdminEmail          string
	AdminPass           string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgresql://core:core@postgres:5432/core?sslmode=disable"),
		ServerAddr:          getEnv("SERVER_ADDR", ":6969"),
		Env:                 getEnv("ENV", "development"),
		HealthCheckInterval: 60 * time.Second,
		SessionSecret:       getEnv("SESSION_SECRET", ""),
		AllowedOrigins:      getEnv("ALLOWED_ORIGINS", "*"),
		AdminEmail:          os.Getenv("ADMIN_EMAIL"),
		AdminPass:           os.Getenv("ADMIN_PASSWORD"),
	}
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

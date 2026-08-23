package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds application configuration
type Config struct {
	DatabaseURL         string
	Port                int
	Env                 string
	FrontendURL         string
	HealthCheckInterval time.Duration
	AdminEmail          string
	AdminPass           string
	AllowedOrigins      string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgresql://core:core@localhost:5432/core?sslmode=disable"),
		Port:                getEnvInt("PORT", 8080),
		Env:                 getEnv("ENV", "development"),
		HealthCheckInterval: 60 * time.Second,
		AdminEmail:          os.Getenv("ADMIN_EMAIL"),
		AdminPass:           os.Getenv("ADMIN_PASSWORD"),
		AllowedOrigins:      getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
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

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return defaultValue
}

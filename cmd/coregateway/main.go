package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ejsadiarin/coregateway/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg := server.Config{
		Port:                getEnvOrDefaultInt("PORT", 8080),
		DatabaseURL:         getEnvOrDefault("DATABASE_URL", "postgresql://core:core@localhost:5432/core?sslmode=disable"),
		FrontendURL:         os.Getenv("FRONTEND_URL"),
		HealthCheckInterval: 60 * time.Second,
		AdminEmail:          os.Getenv("ADMIN_EMAIL"),
		AdminPassword:       os.Getenv("ADMIN_PASSWORD"),
	}

	srv, err := server.New(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	logger.Info("Starting coregateway...", "port", cfg.Port)
	if err := srv.Start(); err != nil {
		logger.Error("Server error", "error", err)
	}

	<-done
	logger.Info("Graceful shutdown complete")
}

func gracefulShutdown(srv *server.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var n int
		if _, err := fmt.Sscanf(value, "%d", &n); err == nil {
			return n
		}
	}
	return defaultValue
}

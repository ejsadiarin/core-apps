package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ejsadiarin/coregateway/internal/config"
	dbpkg "github.com/ejsadiarin/coregateway/internal/db"
	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"
	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	monitor "github.com/ejsadiarin/coregateway/internal/domain/monitor"
	"github.com/ejsadiarin/coregateway/internal/domain/user"
	"github.com/ejsadiarin/coregateway/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	pool, err := dbpkg.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)

	healthChecker := monitor.NewHealthChecker(queries, logger)
	if cfg.HealthCheckInterval > 0 {
		healthChecker.StartHealthCheckScheduler(cfg.HealthCheckInterval)
	}

	authSvc := auth.NewService(queries, logger)
	userSvc := user.NewService(queries, logger)
	monitorSvc := monitor.NewService(queries, logger)

	srv := server.New(cfg, logger, authSvc, userSvc, monitorSvc)

	go startSessionCleanup(queries, logger)

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

func startSessionCleanup(queries *sqlc.Queries, logger *slog.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	logger.Info("Session cleanup scheduler started")

	for range ticker.C {
		if err := queries.DeleteExpiredSessions(context.Background()); err != nil {
			logger.Error("Failed to delete expired sessions", "error", err)
		} else {
			logger.Debug("Expired sessions cleaned up")
		}
	}
}

package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	monitor "github.com/ejsadiarin/coregateway/internal/domain/monitor"
	"github.com/ejsadiarin/coregateway/internal/domain/user"
	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// Server holds all the application dependencies
type Server struct {
	port int

	Echo    *echo.Echo
	DB      *pgxpool.Pool
	Queries *sqlc.Queries
	Logger  *slog.Logger

	HealthChecker  *monitor.HealthChecker
	AuthHandler    *auth.Handler
	UserHandler    *user.Handler
	ServiceHandler *monitor.Handler
}

// Config holds server configuration
type Config struct {
	Port               int
	DatabaseURL        string
	AdminEmail         string
	AdminPassword      string
	HealthCheckInterval time.Duration
	FrontendURL        string
}

// New creates a new Server instance with all dependencies initialized
func New(cfg Config, logger *slog.Logger) (*Server, error) {
	// initialize database connection pool
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	// test the connection
	if err = dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established")

	// initialize sqlc queries
	queries := sqlc.New(dbPool)

	// seed initial users (admin and demo)
	if err := auth.SeedUsers(context.Background(), queries, logger, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		logger.Warn("Failed to seed users", "error", err)
	}

	// initialize health checker
	healthChecker := monitor.NewHealthChecker(queries, logger)
	if cfg.HealthCheckInterval > 0 {
		healthChecker.StartHealthCheckScheduler(cfg.HealthCheckInterval)
	}

	// initialize handlers
	authHandler := auth.NewHandler(queries, logger)
	userHandler := user.NewHandler(queries, logger)
	serviceHandler := monitor.NewHandler(queries, logger)

	// initialize Echo
	e := echo.New()
	e.HideBanner = true

	s := &Server{
		port:           cfg.Port,
		Echo:           e,
		DB:             dbPool,
		Queries:        queries,
		Logger:         logger,
		HealthChecker:  healthChecker,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		ServiceHandler: serviceHandler,
	}

	// setup middleware and routes
	s.setupMiddleware(cfg)
	s.RegisterRoutes()

	// start background jobs
	go s.startSessionCleanupScheduler()

	return s, nil
}

// Close cleans up application resources
func (s *Server) Close() {
	if s.DB != nil {
		s.DB.Close()
	}
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	s.Logger.Info("Server starting", "port", port)
	return s.Echo.Start(":" + port)
}

// startSessionCleanupScheduler runs a background goroutine that periodically deletes expired sessions
func (s *Server) startSessionCleanupScheduler() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	s.Logger.Info("Session cleanup scheduler started")

	for range ticker.C {
		if err := s.Queries.DeleteExpiredSessions(context.Background()); err != nil {
			s.Logger.Error("Failed to delete expired sessions", "error", err)
		} else {
			s.Logger.Debug("Expired sessions cleaned up")
		}
	}
}

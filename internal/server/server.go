package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"
	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	monitor "github.com/ejsadiarin/coregateway/internal/domain/monitor"
	"github.com/ejsadiarin/coregateway/internal/domain/user"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server holds all the application dependencies
type Server struct {
	port int

	router *chi.Mux
	http   *http.Server
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
	Port                int
	DatabaseURL         string
	AdminEmail          string
	AdminPassword       string
	HealthCheckInterval time.Duration
	FrontendURL         string
}

// New creates a new Server instance with all dependencies initialized
func New(cfg Config, logger *slog.Logger) (*Server, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	if err = dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established")

	queries := sqlc.New(dbPool)

	if err := auth.SeedUsers(context.Background(), queries, logger, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		logger.Warn("Failed to seed users", "error", err)
	}

	healthChecker := monitor.NewHealthChecker(queries, logger)
	if cfg.HealthCheckInterval > 0 {
		healthChecker.StartHealthCheckScheduler(cfg.HealthCheckInterval)
	}

	authHandler := auth.NewHandler(queries, logger)
	userHandler := user.NewHandler(queries, logger)
	serviceHandler := monitor.NewHandler(queries, logger)

	router := chi.NewRouter()

	s := &Server{
		port:           cfg.Port,
		router:         router,
		DB:             dbPool,
		Queries:        queries,
		Logger:         logger,
		HealthChecker:  healthChecker,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		ServiceHandler: serviceHandler,
	}

	s.setupMiddleware(cfg)
	s.RegisterRoutes()

	s.http = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

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
func (s *Server) Start() error {
	s.Logger.Info("Server starting", "addr", s.http.Addr)
	return s.http.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
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

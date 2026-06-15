package app

import (
	"context"
	"time"

	"core-gateway/internal/domain/auth"
	"core-gateway/internal/domain/budget"
	"core-gateway/internal/domain/services"
	"core-gateway/internal/domain/user"
	"core-gateway/internal/repository/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// Application holds all the application dependencies
type Application struct {
	Echo          *echo.Echo
	DB            *pgxpool.Pool
	Queries       *sqlc.Queries
	Logger        *zerolog.Logger
	HealthChecker *services.HealthChecker

	// handlers
	AuthHandler    *auth.Handler
	UserHandler    *user.Handler
	ServiceHandler *services.Handler
	BudgetHandler  *budget.Handler
}

// Config holds application configuration
type Config struct {
	DatabaseURL         string
	Port                string
	Env                 string
	FrontendURL         string
	HealthCheckInterval time.Duration
	AdminEmail          string
	AdminPassword       string
}

// New creates a new Application instance with all dependencies initialized
func New(cfg Config, logger *zerolog.Logger) (*Application, error) {
	// initialize database connection pool
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// test the connection
	if err = dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		return nil, err
	}

	logger.Info().Msg("Database connection established")

	// initialize sqlc queries
	queries := sqlc.New(dbPool)

	// seed initial users (admin and demo)
	if err := auth.SeedUsers(context.Background(), queries, logger, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		logger.Warn().Err(err).Msg("Failed to seed users")
	}

	// initialize health checker
	healthChecker := services.NewHealthChecker(queries, logger)
	if cfg.HealthCheckInterval > 0 {
		healthChecker.StartHealthCheckScheduler(cfg.HealthCheckInterval)
	}

	// initialize handlers
	authHandler := auth.NewHandler(queries, logger)
	userHandler := user.NewHandler(queries, logger)
	serviceHandler := services.NewHandler(queries, logger)
	budgetHandler := budget.NewHandler(queries, logger, dbPool)

	// initialize Echo
	e := echo.New()
	e.HideBanner = true

	app := &Application{
		Echo:           e,
		DB:             dbPool,
		Queries:        queries,
		Logger:         logger,
		HealthChecker:  healthChecker,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		ServiceHandler: serviceHandler,
		BudgetHandler:  budgetHandler,
	}

	// setup middleware and routes
	app.SetupMiddleware(cfg)
	app.RegisterRoutes()

	// start background jobs
	go app.startSessionCleanupScheduler()

	return app, nil
}

// Close cleans up application resources
func (a *Application) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}

// Start starts the HTTP server
func (a *Application) Start(port string) error {
	a.Logger.Info().Str("port", port).Msg("Server starting")
	return a.Echo.Start(":" + port)
}

// startSessionCleanupScheduler runs a background goroutine that periodically deletes expired sessions
func (a *Application) startSessionCleanupScheduler() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	a.Logger.Info().Msg("Session cleanup scheduler started")

	for range ticker.C {
		if err := a.Queries.DeleteExpiredSessions(context.Background()); err != nil {
			a.Logger.Error().Err(err).Msg("Failed to delete expired sessions")
		} else {
			a.Logger.Debug().Msg("Expired sessions cleaned up")
		}
	}
}

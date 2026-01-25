// @title Core Homelab API
// @version 2.0
// @description API for Core Homelab Dashboard with service monitoring and budget tracking
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

package main

import (
	"context"
	"core-gateway/handlers"
	"core-gateway/health"
	"core-gateway/internal/middleware"
	"core-gateway/internal/sqlc"
	customValidator "core-gateway/internal/validator"
	"core-gateway/models"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// SystemStats mimics the data structure expected by the frontend
type SystemStats struct {
	CPU         int    `json:"cpu"`
	Memory      int    `json:"memory"`
	Storage     int    `json:"storage"`
	Temperature int    `json:"temperature"`
	Uptime      string `json:"uptime"`
	Network     struct {
		Up   string `json:"up"`
		Down string `json:"down"`
	} `json:"network"`
}

// ServiceStatus represents a single service's state (legacy)
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

func main() {
	// Initialize zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if os.Getenv("ENV") == "production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	logger.Info().Msg("Starting Core Homelab API")

	// Initialize database connection pool
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://core:core@postgres:5432/core?sslmode=disable"
	}

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to parse database URL")
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer dbPool.Close()

	// Test the connection
	if err = dbPool.Ping(context.Background()); err != nil {
		logger.Fatal().Err(err).Msg("Failed to ping database")
	}

	logger.Info().Msg("Database connection established")

	// Initialize sqlc queries
	queries := sqlc.New(dbPool)

	// Initialize health checker
	checker := health.NewChecker(queries, &logger)
	checker.StartHealthCheckScheduler(60 * time.Second)

	// Initialize handlers
	serviceHandler := handlers.NewServiceHandler(queries, &logger)

	// Initialize Echo
	e := echo.New()
	e.HideBanner = true

	// Custom validator
	e.Validator = customValidator.NewValidator()

	// Middleware
	e.Use(middleware.RequestIDMiddleware)
	e.Use(middleware.ZerologMiddleware(middleware.ZerologConfig{
		Logger: &logger,
		Skipper: func(c echo.Context) bool {
			// Skip logging for health check and docs
			return c.Path() == "/health" || c.Path() == "/swagger/*"
		},
	}))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:3001", os.Getenv("FRONTEND_URL")},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))

	// Custom error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprintf("%v", he.Message)
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				c.NoContent(code)
			} else {
				c.JSON(code, models.ErrorResponse{
					Error: message,
				})
			}
		}

		logger.Error().
			Err(err).
			Int("status", code).
			Str("path", c.Request().URL.Path).
			Msg("Request error")
	}

	// Routes
	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Health endpoint
	e.GET("/health", healthCheck)

	// Legacy endpoints (for backwards compatibility)
	e.GET("/api/system/stats", getSystemStats)
	e.GET("/api/services", getLegacyServices)

	// API v1 - Service monitoring
	api := e.Group("/api")
	{
		services := api.Group("/services")
		{
			services.POST("", serviceHandler.CreateService)
			services.GET("/list", serviceHandler.ListServices)
			services.GET("/:id", serviceHandler.GetService)
			services.PUT("/:id", serviceHandler.UpdateService)
			services.DELETE("/:id", serviceHandler.DeleteService)
			services.GET("/:id/history", serviceHandler.GetServiceHistory)
			services.GET("/:id/stats", serviceHandler.GetServiceStats)
			services.GET("/stats/all", serviceHandler.GetAllServicesStats)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info().Str("port", port).Msg("Server starting")
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
}

// healthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// getSystemStats godoc
// @Summary Get system statistics
// @Description Get mock system statistics (CPU, memory, etc.)
// @Tags system
// @Produce json
// @Success 200 {object} SystemStats
// @Router /api/system/stats [get]
func getSystemStats(c echo.Context) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	stats := SystemStats{
		CPU:         rand.Intn(30) + 10,
		Memory:      rand.Intn(40) + 20,
		Storage:     68,
		Temperature: 45,
		Uptime:      "42d 13h 27m",
		Network: struct {
			Up   string `json:"up"`
			Down string `json:"down"`
		}{
			Up:   "125.4 Mbps",
			Down: "342.8 Mbps",
		},
	}

	return c.JSON(http.StatusOK, stats)
}

// getLegacyServices godoc
// @Summary Get legacy services (deprecated)
// @Description Get mock service list for backwards compatibility
// @Tags legacy
// @Produce json
// @Success 200 {array} ServiceStatus
// @Deprecated true
// @Router /api/services [get]
func getLegacyServices(c echo.Context) error {
	services := []ServiceStatus{
		{Name: "Docker Manager", Type: "Container", Status: "online"},
		{Name: "PostgreSQL", Type: "Database", Status: "online"},
		{Name: "Nextcloud", Type: "Storage", Status: "online"},
		{Name: "Vault", Type: "Security", Status: "online"},
		{Name: "Jellyfin", Type: "Media", Status: "offline"},
		{Name: "Mail Server", Type: "Email", Status: "maintenance"},
	}

	return c.JSON(http.StatusOK, services)
}

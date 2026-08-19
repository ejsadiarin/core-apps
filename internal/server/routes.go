package server

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/ejsadiarin/coregateway/internal/domain/auth"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// legacy types for backwards compatibility
type systemStats struct {
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

type serviceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// RegisterRoutes sets up all API routes
func (s *Server) RegisterRoutes() {
	// swagger documentation
	s.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// health endpoint
	s.Echo.GET("/health", s.healthCheck)

	// legacy endpoints (for backwards compatibility)
	s.Echo.GET("/api/system/stats", s.getSystemStats)
	s.Echo.GET("/api/services", s.getLegacyServices)

	// API v1
	api := s.Echo.Group("/api")
	{
		// auth routes (public)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", s.AuthHandler.Register)
			authGroup.POST("/login", s.AuthHandler.Login)
			authGroup.POST("/logout", s.AuthHandler.Logout)
			authGroup.POST("/demo", s.AuthHandler.LoginAsDemo)
			authGroup.GET("/me", s.AuthHandler.Me)
		}

		// user management routes (admin only, except get/update own)
		users := api.Group("/users")
		users.Use(auth.RequireAuth())
		{
			users.GET("", s.UserHandler.ListUsers, auth.RequireRole(auth.RoleAdmin))
			users.POST("", s.UserHandler.CreateUser, auth.RequireRole(auth.RoleAdmin))
			users.GET("/:id", s.UserHandler.GetUser)
			users.PUT("/:id", s.UserHandler.UpdateUser)
			users.DELETE("/:id", s.UserHandler.DeleteUser, auth.RequireRole(auth.RoleAdmin))
		}

		// service monitoring routes
		services := api.Group("/services")
		{
			services.POST("", s.ServiceHandler.CreateService)
			services.GET("/list", s.ServiceHandler.ListServices)
			services.GET("/:id", s.ServiceHandler.GetService)
			services.PUT("/:id", s.ServiceHandler.UpdateService)
			services.DELETE("/:id", s.ServiceHandler.DeleteService)
			services.GET("/:id/history", s.ServiceHandler.GetServiceHistory)
			services.GET("/:id/stats", s.ServiceHandler.GetServiceStats)
			services.GET("/stats/all", s.ServiceHandler.GetAllServicesStats)
		}
	}
}

// healthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (s *Server) healthCheck(c echo.Context) error {
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
// @Success 200 {object} systemStats
// @Router /api/system/stats [get]
func (s *Server) getSystemStats(c echo.Context) error {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	stats := systemStats{
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
// @Success 200 {array} serviceStatus
// @Deprecated true
// @Router /api/services [get]
func (s *Server) getLegacyServices(c echo.Context) error {
	services := []serviceStatus{
		{Name: "Docker Manager", Type: "Container", Status: "online"},
		{Name: "PostgreSQL", Type: "Database", Status: "online"},
		{Name: "Nextcloud", Type: "Storage", Status: "online"},
		{Name: "Vault", Type: "Security", Status: "online"},
		{Name: "Jellyfin", Type: "Media", Status: "offline"},
		{Name: "Mail Server", Type: "Email", Status: "maintenance"},
	}

	return c.JSON(http.StatusOK, services)
}

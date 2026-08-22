package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	"github.com/go-chi/chi/v5"
	swagger "github.com/swaggo/http-swagger"
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
	s.router.Get("/swagger/*", swagger.WrapHandler)

	// health endpoint
	s.router.Get("/health", s.healthCheck)

	// legacy endpoints
	s.router.Get("/api/system/stats", s.getSystemStats)
	s.router.Get("/api/legacy/services", s.getLegacyServices)

	// API v1
	s.router.Route("/api", func(r chi.Router) {
		// auth routes (public)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", s.AuthHandler.Register)
			r.Post("/login", s.AuthHandler.Login)
			r.Post("/logout", s.AuthHandler.Logout)
			r.Post("/demo", s.AuthHandler.LoginAsDemo)
			r.Get("/me", s.AuthHandler.Me)
		})

		// user management routes (require auth)
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth())

			r.Get("/users", s.UserHandler.ListUsers)
			r.Post("/users", s.UserHandler.CreateUser)
			r.Get("/users/{id}", s.UserHandler.GetUser)
			r.Put("/users/{id}", s.UserHandler.UpdateUser)

			// admin-only routes
			r.Group(func(r chi.Router) {
				r.Use(auth.RequireRole(auth.RoleAdmin))
				r.Delete("/users/{id}", s.UserHandler.DeleteUser)
			})
		})

		// service monitoring routes
		r.Route("/services", func(r chi.Router) {
			r.Post("", s.ServiceHandler.CreateService)
			r.Get("/list", s.ServiceHandler.ListServices)
			r.Get("/stats/all", s.ServiceHandler.GetAllServicesStats)
			r.Get("/{id}", s.ServiceHandler.GetService)
			r.Put("/{id}", s.ServiceHandler.UpdateService)
			r.Delete("/{id}", s.ServiceHandler.DeleteService)
			r.Get("/{id}/history", s.ServiceHandler.GetServiceHistory)
			r.Get("/{id}/stats", s.ServiceHandler.GetServiceStats)
		})
	})
}

// healthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
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
func (s *Server) getSystemStats(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// getLegacyServices godoc
// @Summary Get legacy services (deprecated)
// @Description Get mock service list for backwards compatibility
// @Tags legacy
// @Produce json
// @Success 200 {array} serviceStatus
// @Deprecated true
// @Router /api/legacy/services [get]
func (s *Server) getLegacyServices(w http.ResponseWriter, r *http.Request) {
	services := []serviceStatus{
		{Name: "Docker Manager", Type: "Container", Status: "online"},
		{Name: "PostgreSQL", Type: "Database", Status: "online"},
		{Name: "Nextcloud", Type: "Storage", Status: "online"},
		{Name: "Vault", Type: "Security", Status: "online"},
		{Name: "Jellyfin", Type: "Media", Status: "offline"},
		{Name: "Mail Server", Type: "Email", Status: "maintenance"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

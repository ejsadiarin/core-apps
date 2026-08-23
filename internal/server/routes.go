package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	swagger "github.com/swaggo/http-swagger"
)

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

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// global middleware
	r.Use(auth.AuthMiddleware(s.authService))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", s.frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// ops
	r.Get("/swagger/*", swagger.WrapHandler)
	r.Get("/health", healthCheck)
	r.Get("/api/system/stats", getSystemStats)
	r.Get("/api/legacy/services", getLegacyServices)

	// auth
	r.Post("/api/auth/register", s.authHandler.Register)
	r.Post("/api/auth/login", s.authHandler.Login)
	r.Post("/api/auth/logout", s.authHandler.Logout)
	r.Post("/api/auth/demo", s.authHandler.LoginAsDemo)
	r.Get("/api/auth/me", s.authHandler.Me)

	// users
	r.Get("/api/users", s.userHandler.ListUsers)
	r.Post("/api/users", s.userHandler.CreateUser)
	r.Get("/api/users/{id}", s.userHandler.GetUser)
	r.Put("/api/users/{id}", s.userHandler.UpdateUser)
	r.Delete("/api/users/{id}", s.userHandler.DeleteUser)

	// services
	r.Post("/api/services", s.serviceHandler.CreateService)
	r.Get("/api/services/list", s.serviceHandler.ListServices)
	r.Get("/api/services/stats/all", s.serviceHandler.GetAllServicesStats)
	r.Get("/api/services/{id}", s.serviceHandler.GetService)
	r.Put("/api/services/{id}", s.serviceHandler.UpdateService)
	r.Delete("/api/services/{id}", s.serviceHandler.DeleteService)
	r.Get("/api/services/{id}/history", s.serviceHandler.GetServiceHistory)
	r.Get("/api/services/{id}/stats", s.serviceHandler.GetServiceStats)

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func getSystemStats(w http.ResponseWriter, r *http.Request) {
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

func getLegacyServices(w http.ResponseWriter, r *http.Request) {
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

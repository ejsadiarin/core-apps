package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

// ServiceStatus represents a single service's state
type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"` // online, offline, maintenance
	Type   string `json:"type"`
}

func main() {
	// Initialize Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // Allow frontend
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	e.GET("/health", healthCheck)
	e.GET("/api/system/stats", getSystemStats)
	e.GET("/api/services", getServices)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// healthCheck returns a simple 200 OK
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// getSystemStats returns mock data for the dashboard charts/circles
// TODO: Replace with Prometheus queries
func getSystemStats(c echo.Context) error {
	// Mock randomization to simulate live data
	rand.New(rand.NewSource(time.Now().UnixNano()))

	stats := SystemStats{
		CPU:         rand.Intn(30) + 10, // 10-40%
		Memory:      rand.Intn(40) + 20, // 20-60%
		Storage:     68,                 // Static for now
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

// getServices returns mock service status
// TODO: Replace with K8s API calls
func getServices(c echo.Context) error {
	services := []ServiceStatus{
		{Name: "Docker Manager", Type: "Container", Status: "online"},
		{Name: "PostgreSQL", Type: "Database", Status: "online"},
		{Name: "Nextcloud", Type: "Storage", Status: "online"},
		{Name: "Vault", Type: "Security", Status: "online"},
		{Name: "Jellyfin", Type: "Media", Status: "online"},
		{Name: "Mail Server", Type: "Email", Status: "maintenance"},
	}

	return c.JSON(http.StatusOK, services)
}

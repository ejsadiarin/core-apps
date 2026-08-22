package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

type Handler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

func NewHandler(queries *sqlc.Queries, logger *slog.Logger) *Handler {
	return &Handler{
		queries: queries,
		logger:  logger,
	}
}

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service for health monitoring
// @Tags services
// @Accept json
// @Produce json
// @Param service body CreateServiceRequest true "Service to create"
// @Success 201 {object} sqlc.Service
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/services [post]
func (h *Handler) CreateService(w http.ResponseWriter, r *http.Request) {
	var req CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	healthCheckInterval := int32(60)
	if req.HealthCheckInterval != nil {
		healthCheckInterval = *req.HealthCheckInterval
	}

	healthCheckMethod := "GET"
	if req.HealthCheckMethod != nil {
		healthCheckMethod = *req.HealthCheckMethod
	}

	timeout := int32(5000)
	if req.Timeout != nil {
		timeout = *req.Timeout
	}

	expectedStatusCodes := []int32{200, 204}
	if len(req.ExpectedStatusCodes) > 0 {
		expectedStatusCodes = req.ExpectedStatusCodes
	}

	var icon, description, serviceType pgtype.Text
	if req.Icon != nil {
		icon = pgtype.Text{String: *req.Icon, Valid: true}
	}
	if req.Description != nil {
		description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.ServiceType != nil {
		serviceType = pgtype.Text{String: *req.ServiceType, Valid: true}
	}

	service, err := h.queries.CreateService(context.Background(), sqlc.CreateServiceParams{
		Name:                req.Name,
		Url:                 req.URL,
		Icon:                icon,
		Description:         description,
		ServiceType:         serviceType,
		HealthCheckInterval: pgtype.Int4{Int32: healthCheckInterval, Valid: true},
		HealthCheckMethod:   pgtype.Text{String: healthCheckMethod, Valid: true},
		ExpectedStatusCodes: expectedStatusCodes,
		Timeout:             pgtype.Int4{Int32: timeout, Valid: true},
	})
	if err != nil {
		h.logger.Error("Failed to create service", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create service"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(service)
}

// ListServices godoc
// @Summary List all services
// @Description Get all active services with their current health status
// @Tags services
// @Produce json
// @Success 200 {array} sqlc.ListServicesRow
// @Failure 500 {object} ErrorResponse
// @Router /api/services/list [get]
func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
	services, err := h.queries.ListServices(context.Background())
	if err != nil {
		h.logger.Error("Failed to list services", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to list services"})
		return
	}

	if services == nil {
		services = []sqlc.ListServicesRow{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

// GetService godoc
// @Summary Get a service by ID
// @Description Get detailed information about a specific service
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} sqlc.Service
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/services/{id} [get]
func (h *Handler) GetService(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid service ID"})
		return
	}

	service, err := h.queries.GetService(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Service not found"})
			return
		}
		h.logger.Error("Failed to get service", "error", err, "id", id.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get service"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

// UpdateService godoc
// @Summary Update a service
// @Description Update an existing service configuration
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body UpdateServiceRequest true "Service updates"
// @Success 200 {object} sqlc.Service
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/services/{id} [put]
func (h *Handler) UpdateService(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid service ID"})
		return
	}

	var req UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	params := sqlc.UpdateServiceParams{ID: id}

	if req.Name != nil {
		params.Name = pgtype.Text{String: *req.Name, Valid: true}
	}
	if req.URL != nil {
		params.Url = pgtype.Text{String: *req.URL, Valid: true}
	}
	if req.Icon != nil {
		params.Icon = pgtype.Text{String: *req.Icon, Valid: true}
	}
	if req.Description != nil {
		params.Description = pgtype.Text{String: *req.Description, Valid: true}
	}
	if req.ServiceType != nil {
		params.ServiceType = pgtype.Text{String: *req.ServiceType, Valid: true}
	}
	if req.HealthCheckInterval != nil {
		params.HealthCheckInterval = pgtype.Int4{Int32: *req.HealthCheckInterval, Valid: true}
	}
	if req.HealthCheckMethod != nil {
		params.HealthCheckMethod = pgtype.Text{String: *req.HealthCheckMethod, Valid: true}
	}
	if req.ExpectedStatusCodes != nil {
		params.ExpectedStatusCodes = req.ExpectedStatusCodes
	}
	if req.Timeout != nil {
		params.Timeout = pgtype.Int4{Int32: *req.Timeout, Valid: true}
	}
	if req.IsActive != nil {
		params.IsActive = pgtype.Bool{Bool: *req.IsActive, Valid: true}
	}

	service, err := h.queries.UpdateService(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Service not found"})
			return
		}
		h.logger.Error("Failed to update service", "error", err, "id", id.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to update service"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(service)
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete a service by ID
// @Tags services
// @Param id path string true "Service ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/services/{id} [delete]
func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid service ID"})
		return
	}

	err = h.queries.DeleteService(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Service not found"})
			return
		}
		h.logger.Error("Failed to delete service", "error", err, "id", id.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to delete service"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetServiceHistory godoc
// @Summary Get service health history
// @Description Get health check history for a specific service
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Param limit query int false "Limit" default(100)
// @Success 200 {array} sqlc.ServiceHealthHistory
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/services/{id}/history [get]
func (h *Handler) GetServiceHistory(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid service ID"})
		return
	}

	limit := int32(100)

	history, err := h.queries.GetServiceHistory(context.Background(), sqlc.GetServiceHistoryParams{
		ServiceID: pgtype.UUID{Bytes: id, Valid: true},
		Limit:     limit,
	})
	if err != nil {
		h.logger.Error("Failed to get service history", "error", err, "id", id.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get service history"})
		return
	}

	if history == nil {
		history = []sqlc.CoregatewayServiceHealthHistory{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// GetServiceStats godoc
// @Summary Get service statistics
// @Description Get uptime statistics for a specific service
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/services/{id}/stats [get]
func (h *Handler) GetServiceStats(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid service ID"})
		return
	}

	pgID := pgtype.UUID{Bytes: id, Valid: true}

	stats24h, err := h.queries.GetServiceStats24h(context.Background(), pgID)
	if err != nil {
		h.logger.Error("Failed to get 24h stats", "error", err, "id", id.String())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get service stats"})
		return
	}

	stats7d, _ := h.queries.GetServiceStats7d(context.Background(), pgID)
	stats30d, _ := h.queries.GetServiceStats30d(context.Background(), pgID)

	uptime24h := 0.0
	if stats24h.TotalChecks > 0 {
		uptime24h = float64(stats24h.SuccessfulChecks) / float64(stats24h.TotalChecks) * 100
	}

	uptime7d := 0.0
	if stats7d.TotalChecks > 0 {
		uptime7d = float64(stats7d.SuccessfulChecks) / float64(stats7d.TotalChecks) * 100
	}

	uptime30d := 0.0
	if stats30d.TotalChecks > 0 {
		uptime30d = float64(stats30d.SuccessfulChecks) / float64(stats30d.TotalChecks) * 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service_id":        id.String(),
		"uptime_24h":        uptime24h,
		"uptime_7d":         uptime7d,
		"uptime_30d":        uptime30d,
		"avg_response_time": stats24h.AvgResponseTime,
		"total_checks":      stats24h.TotalChecks,
		"successful_checks": stats24h.SuccessfulChecks,
	})
}

// GetAllServicesStats godoc
// @Summary Get overall statistics
// @Description Get overall statistics for all services
// @Tags services
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} ErrorResponse
// @Router /api/services/stats/all [get]
func (h *Handler) GetAllServicesStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.queries.GetAllServicesStats(context.Background())
	if err != nil {
		h.logger.Error("Failed to get all services stats", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to get services stats"})
		return
	}

	uptime := 0.0
	if stats.TotalChecks > 0 {
		uptime = float64(stats.SuccessfulChecks) / float64(stats.TotalChecks) * 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_services":    stats.TotalServices,
		"total_checks":      stats.TotalChecks,
		"successful_checks": stats.SuccessfulChecks,
		"overall_uptime":    uptime,
		"avg_response_time": stats.AvgResponseTime,
	})
}

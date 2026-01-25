package handlers

import (
	"context"
	"core-gateway/internal/sqlc"
	"core-gateway/internal/validator"
	"core-gateway/models"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ServiceHandler struct {
	queries *sqlc.Queries
	logger  *zerolog.Logger
}

func NewServiceHandler(queries *sqlc.Queries, logger *zerolog.Logger) *ServiceHandler {
	return &ServiceHandler{
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
// @Param service body models.CreateServiceRequest true "Service to create"
// @Success 201 {object} sqlc.Service
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/services [post]
func (h *ServiceHandler) CreateService(c echo.Context) error {
	req, err := validator.BindAndValidate[models.CreateServiceRequest](c)
	if err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
	}

	// Set defaults
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

	// Convert to pgtype.Text for nullable fields
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
		h.logger.Error().Err(err).Msg("Failed to create service")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create service",
		})
	}

	return c.JSON(http.StatusCreated, service)
}

// ListServices godoc
// @Summary List all services
// @Description Get all active services with their current health status
// @Tags services
// @Produce json
// @Success 200 {array} sqlc.ListServicesRow
// @Failure 500 {object} models.ErrorResponse
// @Router /api/services/list [get]
func (h *ServiceHandler) ListServices(c echo.Context) error {
	services, err := h.queries.ListServices(context.Background())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list services")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to list services",
		})
	}

	if services == nil {
		services = []sqlc.ListServicesRow{}
	}

	return c.JSON(http.StatusOK, services)
}

// GetService godoc
// @Summary Get a service by ID
// @Description Get detailed information about a specific service
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} sqlc.Service
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/services/{id} [get]
func (h *ServiceHandler) GetService(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid service ID",
		})
	}

	service, err := h.queries.GetService(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Service not found",
			})
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("Failed to get service")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get service",
		})
	}

	return c.JSON(http.StatusOK, service)
}

// UpdateService godoc
// @Summary Update a service
// @Description Update an existing service configuration
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body models.UpdateServiceRequest true "Service updates"
// @Success 200 {object} sqlc.Service
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/services/{id} [put]
func (h *ServiceHandler) UpdateService(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid service ID",
		})
	}

	req, err := validator.BindAndValidate[models.UpdateServiceRequest](c)
	if err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
	}

	// Build params with nullable types
	params := sqlc.UpdateServiceParams{
		ID: id,
	}

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
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Service not found",
			})
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("Failed to update service")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update service",
		})
	}

	return c.JSON(http.StatusOK, service)
}

// DeleteService godoc
// @Summary Delete a service
// @Description Delete a service by ID
// @Tags services
// @Param id path string true "Service ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/services/{id} [delete]
func (h *ServiceHandler) DeleteService(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid service ID",
		})
	}

	err = h.queries.DeleteService(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Service not found",
			})
		}
		h.logger.Error().Err(err).Str("id", id.String()).Msg("Failed to delete service")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete service",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetServiceHistory godoc
// @Summary Get service health history
// @Description Get health check history for a specific service
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Param limit query int false "Limit" default(100)
// @Success 200 {array} sqlc.ServiceHealthHistory
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/services/{id}/history [get]
func (h *ServiceHandler) GetServiceHistory(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid service ID",
		})
	}

	limit := int32(100)

	history, err := h.queries.GetServiceHistory(context.Background(), sqlc.GetServiceHistoryParams{
		ServiceID: pgtype.UUID{Bytes: id, Valid: true},
		Limit:     limit,
	})
	if err != nil {
		h.logger.Error().Err(err).Str("id", id.String()).Msg("Failed to get service history")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get service history",
		})
	}

	if history == nil {
		history = []sqlc.ServiceHealthHistory{}
	}

	return c.JSON(http.StatusOK, history)
}

// GetServiceStats godoc
// @Summary Get service statistics
// @Description Get uptime statistics for a specific service
// @Tags services
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/services/{id}/stats [get]
func (h *ServiceHandler) GetServiceStats(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid service ID",
		})
	}

	pgID := pgtype.UUID{Bytes: id, Valid: true}

	stats24h, err := h.queries.GetServiceStats24h(context.Background(), pgID)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id.String()).Msg("Failed to get 24h stats")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get service stats",
		})
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

	return c.JSON(http.StatusOK, map[string]interface{}{
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
// @Failure 500 {object} models.ErrorResponse
// @Router /api/services/stats/all [get]
func (h *ServiceHandler) GetAllServicesStats(c echo.Context) error {
	stats, err := h.queries.GetAllServicesStats(context.Background())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get all services stats")
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to get services stats",
		})
	}

	uptime := 0.0
	if stats.TotalChecks > 0 {
		uptime = float64(stats.SuccessfulChecks) / float64(stats.TotalChecks) * 100
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total_services":    stats.TotalServices,
		"total_checks":      stats.TotalChecks,
		"successful_checks": stats.SuccessfulChecks,
		"overall_uptime":    uptime,
		"avg_response_time": stats.AvgResponseTime,
	})
}

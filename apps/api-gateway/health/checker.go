package health

import (
	"context"
	"core-gateway/internal/sqlc"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

type Checker struct {
	queries *sqlc.Queries
	logger  *zerolog.Logger
}

func NewChecker(queries *sqlc.Queries, logger *zerolog.Logger) *Checker {
	return &Checker{
		queries: queries,
		logger:  logger,
	}
}

// CheckService performs a health check on a single service
func (ch *Checker) CheckService(service sqlc.ListActiveServicesForHealthCheckRow) {
	start := time.Now()

	// Extract health check method
	method := "GET"
	if service.HealthCheckMethod.Valid {
		method = service.HealthCheckMethod.String
	}

	// Extract timeout
	timeout := 5000
	if service.Timeout.Valid {
		timeout = int(service.Timeout.Int32)
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	req, err := http.NewRequest(method, service.Url, nil)
	if err != nil {
		ch.saveHealthCheck(service.ID, "offline", 0, 0, fmt.Sprintf("Failed to create request: %v", err))
		return
	}

	resp, err := client.Do(req)
	responseTime := int32(time.Since(start).Milliseconds())

	if err != nil {
		ch.saveHealthCheck(service.ID, "offline", responseTime, 0, fmt.Sprintf("Request failed: %v", err))
		return
	}
	defer resp.Body.Close()

	// Check if status code is expected
	statusOK := false
	for _, code := range service.ExpectedStatusCodes {
		if int32(resp.StatusCode) == code {
			statusOK = true
			break
		}
	}

	status := "online"
	errorMsg := ""
	if !statusOK {
		status = "offline"
		errorMsg = fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)
	} else if responseTime > 5000 {
		status = "degraded"
	}

	ch.saveHealthCheck(service.ID, status, responseTime, int32(resp.StatusCode), errorMsg)
}

// saveHealthCheck saves a health check result to the database
func (ch *Checker) saveHealthCheck(serviceID uuid.UUID, status string, responseTime int32, statusCode int32, errorMsg string) {
	params := sqlc.CreateHealthHistoryParams{
		ServiceID: pgtype.UUID{Bytes: serviceID, Valid: true},
		Status:    status,
	}

	if responseTime > 0 {
		params.ResponseTime = pgtype.Int4{Int32: responseTime, Valid: true}
	}

	if statusCode > 0 {
		params.StatusCode = pgtype.Int4{Int32: statusCode, Valid: true}
	}

	if errorMsg != "" {
		params.ErrorMessage = pgtype.Text{String: errorMsg, Valid: true}
	}

	_, err := ch.queries.CreateHealthHistory(context.Background(), params)
	if err != nil {
		ch.logger.Error().Err(err).Msg("Failed to save health check")
		return
	}

	ch.logger.Info().
		Str("service_id", serviceID.String()).
		Str("status", status).
		Int32("response_time", responseTime).
		Msg("Health check completed")
}

// CheckAllServices performs health checks on all active services
func (ch *Checker) CheckAllServices() {
	services, err := ch.queries.ListActiveServicesForHealthCheck(context.Background())
	if err != nil {
		ch.logger.Error().Err(err).Msg("Failed to fetch services for health check")
		return
	}

	ch.logger.Info().Int("count", len(services)).Msg("Starting health checks")

	// Perform health checks concurrently
	for _, service := range services {
		go ch.CheckService(service)
	}
}

// StartHealthCheckScheduler starts the background health check scheduler
func (ch *Checker) StartHealthCheckScheduler(interval time.Duration) {
	// Initial check
	go ch.CheckAllServices()

	// Schedule checks
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			ch.CheckAllServices()
		}
	}()

	ch.logger.Info().Dur("interval", interval).Msg("Health check scheduler started")
}

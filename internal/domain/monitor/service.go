package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	sqlc "github.com/ejsadiarin/coregateway/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	CreateService(ctx context.Context, req CreateServiceRequest) (*sqlc.CoregatewayService, error)
	ListServices(ctx context.Context) ([]sqlc.ListServicesRow, error)
	GetService(ctx context.Context, id uuid.UUID) (*sqlc.CoregatewayService, error)
	UpdateService(ctx context.Context, id uuid.UUID, req UpdateServiceRequest) (*sqlc.CoregatewayService, error)
	DeleteService(ctx context.Context, id uuid.UUID) error
	GetServiceHistory(ctx context.Context, serviceID uuid.UUID, limit int32) ([]sqlc.CoregatewayServiceHealthHistory, error)
	GetServiceStats(ctx context.Context, serviceID uuid.UUID) (ServiceStats, error)
	GetAllServiceStats(ctx context.Context) (AllServiceStats, error)
}

type ServiceStats struct {
	ServiceID        string      `json:"service_id"`
	Uptime24h        float64     `json:"uptime_24h"`
	Uptime7d         float64     `json:"uptime_7d"`
	Uptime30d        float64     `json:"uptime_30d"`
	AvgResponseTime  interface{} `json:"avg_response_time"`
	TotalChecks      int64       `json:"total_checks"`
	SuccessfulChecks int64       `json:"successful_checks"`
}

type AllServiceStats struct {
	TotalServices    int64       `json:"total_services"`
	TotalChecks      int64       `json:"total_checks"`
	SuccessfulChecks int64       `json:"successful_checks"`
	OverallUptime    float64     `json:"overall_uptime"`
	AvgResponseTime  interface{} `json:"avg_response_time"`
}

type serviceImpl struct {
	queries sqlc.Querier
	logger  *slog.Logger
}

func NewService(queries sqlc.Querier, logger *slog.Logger) Service {
	return &serviceImpl{queries: queries, logger: logger}
}

func (s *serviceImpl) CreateService(ctx context.Context, req CreateServiceRequest) (*sqlc.CoregatewayService, error) {
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

	service, err := s.queries.CreateService(ctx, sqlc.CreateServiceParams{
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
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return &service, nil
}

func (s *serviceImpl) ListServices(ctx context.Context) ([]sqlc.ListServicesRow, error) {
	services, err := s.queries.ListServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	if services == nil {
		services = []sqlc.ListServicesRow{}
	}

	return services, nil
}

func (s *serviceImpl) GetService(ctx context.Context, id uuid.UUID) (*sqlc.CoregatewayService, error) {
	service, err := s.queries.GetService(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service not found")
		}
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	return &service, nil
}

func (s *serviceImpl) UpdateService(ctx context.Context, id uuid.UUID, req UpdateServiceRequest) (*sqlc.CoregatewayService, error) {
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

	service, err := s.queries.UpdateService(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("service not found")
		}
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	return &service, nil
}

func (s *serviceImpl) DeleteService(ctx context.Context, id uuid.UUID) error {
	err := s.queries.DeleteService(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("service not found")
		}
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

func (s *serviceImpl) GetServiceHistory(ctx context.Context, serviceID uuid.UUID, limit int32) ([]sqlc.CoregatewayServiceHealthHistory, error) {
	history, err := s.queries.GetServiceHistory(ctx, sqlc.GetServiceHistoryParams{
		ServiceID: pgtype.UUID{Bytes: serviceID, Valid: true},
		Limit:     limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get service history: %w", err)
	}

	if history == nil {
		history = []sqlc.CoregatewayServiceHealthHistory{}
	}

	return history, nil
}

func (s *serviceImpl) GetServiceStats(ctx context.Context, serviceID uuid.UUID) (ServiceStats, error) {
	pgID := pgtype.UUID{Bytes: serviceID, Valid: true}

	stats24h, err := s.queries.GetServiceStats24h(ctx, pgID)
	if err != nil {
		return ServiceStats{}, fmt.Errorf("failed to get 24h stats: %w", err)
	}

	stats7d, _ := s.queries.GetServiceStats7d(ctx, pgID)
	stats30d, _ := s.queries.GetServiceStats30d(ctx, pgID)

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

	return ServiceStats{
		ServiceID:        serviceID.String(),
		Uptime24h:        uptime24h,
		Uptime7d:         uptime7d,
		Uptime30d:        uptime30d,
		AvgResponseTime:  stats24h.AvgResponseTime,
		TotalChecks:      stats24h.TotalChecks,
		SuccessfulChecks: stats24h.SuccessfulChecks,
	}, nil
}

func (s *serviceImpl) GetAllServiceStats(ctx context.Context) (AllServiceStats, error) {
	stats, err := s.queries.GetAllServicesStats(ctx)
	if err != nil {
		return AllServiceStats{}, fmt.Errorf("failed to get all services stats: %w", err)
	}

	uptime := 0.0
	if stats.TotalChecks > 0 {
		uptime = float64(stats.SuccessfulChecks) / float64(stats.TotalChecks) * 100
	}

	return AllServiceStats{
		TotalServices:    stats.TotalServices,
		TotalChecks:      stats.TotalChecks,
		SuccessfulChecks: stats.SuccessfulChecks,
		OverallUptime:    uptime,
		AvgResponseTime:  stats.AvgResponseTime,
	}, nil
}

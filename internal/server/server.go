package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ejsadiarin/coregateway/internal/config"
	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	monitor "github.com/ejsadiarin/coregateway/internal/domain/monitor"
	"github.com/ejsadiarin/coregateway/internal/domain/user"
)

type Server struct {
	port           int
	http           *http.Server
	logger         *slog.Logger
	frontendURL    string
	authHandler    *auth.Handler
	userHandler    *user.Handler
	serviceHandler *monitor.Handler
	authService    auth.Service
}

func New(
	cfg *config.Config,
	logger *slog.Logger,
	authSvc auth.Service,
	userSvc user.Service,
	monitorSvc monitor.Service,
) *Server {
	s := &Server{
		port:           cfg.Port,
		logger:         logger,
		frontendURL:    cfg.FrontendURL,
		authHandler:    auth.NewHandler(authSvc),
		userHandler:    user.NewHandler(userSvc),
		serviceHandler: monitor.NewHandler(monitorSvc),
		authService:    authSvc,
	}

	s.http = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      s.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	return s
}

func (s *Server) Start() error {
	s.logger.Info("Server starting", "addr", s.http.Addr)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

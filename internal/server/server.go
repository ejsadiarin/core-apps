package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ejsadiarin/coregateway/internal/auth"
	"github.com/ejsadiarin/coregateway/internal/config"
	monitor "github.com/ejsadiarin/coregateway/internal/monitor"
	"github.com/ejsadiarin/coregateway/internal/services/corefinance"
	"github.com/ejsadiarin/coregateway/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/ejsadiarin/coregateway/internal/db/sqlc"
)

type Server struct {
	port              int
	AuthHandler       *auth.Handler
	UserHandler       *user.Handler
	ServiceHandler    *monitor.Handler
	AuthService       *auth.Service
	CorefinanceClient *corefinance.Client
	Queries           *db.Queries
	Pool              *pgxpool.Pool
}

func New(cfg *config.Config, pool *pgxpool.Pool, queries *db.Queries) *http.Server {
	s := &Server{
		port:    cfg.Port,
		Pool:    pool,
		Queries: queries,
	}

	s.AuthService = auth.NewService(queries)
	s.AuthHandler = auth.NewHandler(s.AuthService)
	s.UserHandler = user.NewHandler(user.NewService(queries))
	s.ServiceHandler = monitor.NewHandler(monitor.NewService(queries))
	s.CorefinanceClient = corefinance.New(cfg.CorefinanceURL)

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      s.RegisterRoutes(cfg),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

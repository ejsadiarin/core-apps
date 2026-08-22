package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ejsadiarin/coregateway/internal/domain/auth"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// loggingMiddleware logs HTTP requests using slog
func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			duration := time.Since(start)

			level := slog.LevelInfo
			if ww.Status() >= 500 {
				level = slog.LevelError
			} else if ww.Status() >= 400 {
				level = slog.LevelWarn
			}

			logger.LogAttrs(r.Context(), level, "Request",
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI),
				slog.Int("status", ww.Status()),
				slog.Duration("duration", duration),
				slog.String("remote_ip", r.RemoteAddr),
			)
		})
	}
}

// setupMiddleware applies all middleware to the router
func (s *Server) setupMiddleware(cfg Config) {
	s.router.Use(middleware.RequestID)
	s.router.Use(loggingMiddleware(s.Logger))
	s.router.Use(middleware.Recoverer)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))
	s.router.Use(auth.AuthMiddleware(s.Queries))
}

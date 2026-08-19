package middleware

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := uuid.New().String()
		c.Request().Header.Set("X-Request-ID", id)
		c.Response().Header().Set("X-Request-ID", id)
		return next(c)
	}
}

// LoggingMiddleware returns a middleware that logs HTTP requests using slog
func LoggingMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start)

			req := c.Request()
			res := c.Response()

			level := slog.LevelInfo
			if res.Status >= 500 {
				level = slog.LevelError
			} else if res.Status >= 400 {
				level = slog.LevelWarn
			}

			logger.LogAttrs(c.Request().Context(), level, "Request",
				slog.String("method", req.Method),
				slog.String("uri", req.RequestURI),
				slog.Int("status", res.Status),
				slog.Duration("duration", duration),
				slog.String("remote_ip", req.RemoteAddr),
			)

			return err
		}
	}
}

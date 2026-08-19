package server

import (
	"fmt"
	"net/http"

	"github.com/ejsadiarin/coregateway/internal/domain/auth"
	"github.com/ejsadiarin/coregateway/internal/shared/middleware"
	"github.com/ejsadiarin/coregateway/internal/shared/models"
	"github.com/ejsadiarin/coregateway/internal/shared/validator"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// setupMiddleware configures all middleware for the application
func (s *Server) setupMiddleware(cfg Config) {
	// custom validator
	s.Echo.Validator = validator.NewValidator()

	// request ID middleware
	s.Echo.Use(middleware.RequestIDMiddleware)

	// logging middleware
	s.Echo.Use(middleware.LoggingMiddleware(s.Logger))

	// recover from panics
	s.Echo.Use(echomiddleware.Recover())

	// CORS configuration
	s.Echo.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", cfg.FrontendURL},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))

	// auth middleware - applies to all requests, extracts user from session if present
	s.Echo.Use(auth.AuthMiddleware(s.Queries))

	// custom error handler
	s.Echo.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = fmt.Sprintf("%v", he.Message)
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				c.NoContent(code)
			} else {
				c.JSON(code, models.ErrorResponse{
					Error: message,
				})
			}
		}

		s.Logger.Error("Request error",
			"error", err,
			"status", code,
			"path", c.Request().URL.Path,
		)
	}
}

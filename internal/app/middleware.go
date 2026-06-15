package app

import (
	"fmt"
	"net/http"

	"core-gateway/internal/domain/auth"
	"core-gateway/internal/shared/middleware"
	"core-gateway/internal/shared/models"
	"core-gateway/internal/shared/validator"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware configures all middleware for the application
func (a *Application) SetupMiddleware(cfg Config) {
	// custom validator
	a.Echo.Validator = validator.NewValidator()

	// request ID middleware
	a.Echo.Use(middleware.RequestIDMiddleware)

	// zerolog logging middleware
	a.Echo.Use(middleware.ZerologMiddleware(middleware.ZerologConfig{
		Logger: a.Logger,
		Skipper: func(c echo.Context) bool {
			// skip logging for health check and docs
			return c.Path() == "/health" || c.Path() == "/swagger/*"
		},
	}))

	// recover from panics
	a.Echo.Use(echomiddleware.Recover())

	// CORS configuration
	a.Echo.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", cfg.FrontendURL},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowCredentials: true,
	}))

	// auth middleware - applies to all requests, extracts user from session if present
	a.Echo.Use(auth.AuthMiddleware(a.Queries))

	// custom error handler
	a.Echo.HTTPErrorHandler = func(err error, c echo.Context) {
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

		a.Logger.Error().
			Err(err).
			Int("status", code).
			Str("path", c.Request().URL.Path).
			Msg("Request error")
	}
}

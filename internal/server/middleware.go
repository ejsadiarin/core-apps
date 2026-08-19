package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ejsadiarin/coregateway/internal/domain/auth"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// CustomValidator wraps go-playground/validator for Echo
type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// BindAndValidate binds the request body to a typed struct and validates it
func BindAndValidate[T any](c echo.Context) (*T, error) {
	var req T
	if err := c.Bind(&req); err != nil {
		return nil, err
	}
	if err := c.Validate(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// FormatValidationErrors formats validator errors into a user-friendly map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors[e.Field()] = e.Tag()
		}
	}
	return errors
}

// requestIDMiddleware adds a unique request ID to each request
func requestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := uuid.New().String()
		c.Request().Header.Set("X-Request-ID", id)
		c.Response().Header().Set("X-Request-ID", id)
		return next(c)
	}
}

// loggingMiddleware returns a middleware that logs HTTP requests using slog
func loggingMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
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

// setupMiddleware configures all middleware for the application
func (s *Server) setupMiddleware(cfg Config) {
	// custom validator
	s.Echo.Validator = NewValidator()

	// request ID middleware
	s.Echo.Use(requestIDMiddleware)

	// logging middleware
	s.Echo.Use(loggingMiddleware(s.Logger))

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
				c.JSON(code, ErrorResponse{
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

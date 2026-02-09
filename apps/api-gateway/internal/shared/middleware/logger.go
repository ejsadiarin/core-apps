package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// ZerologConfig defines the config for Zerolog middleware
type ZerologConfig struct {
	Logger  *zerolog.Logger
	Skipper func(c echo.Context) bool
}

// ZerologMiddleware returns a middleware that logs HTTP requests using zerolog
func ZerologMiddleware(config ZerologConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = func(c echo.Context) bool {
			return false
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()

			var err error
			if err = next(c); err != nil {
				c.Error(err)
			}

			stop := time.Now()
			latency := stop.Sub(start)

			// build log event
			event := config.Logger.Info()

			if err != nil {
				event = config.Logger.Error().Err(err)
			}

			event.
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("remote_ip", c.RealIP()).
				Int("status", res.Status).
				Int64("latency_ms", latency.Milliseconds()).
				Str("user_agent", req.UserAgent()).
				Str("request_id", res.Header().Get(echo.HeaderXRequestID)).
				Int64("bytes_in", req.ContentLength).
				Int64("bytes_out", res.Size).
				Msg("HTTP Request")

			return nil
		}
	}
}

// RequestIDMiddleware adds a request ID to each request if not present
func RequestIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		rid := req.Header.Get(echo.HeaderXRequestID)
		if rid == "" {
			rid = generateRequestID()
			req.Header.Set(echo.HeaderXRequestID, rid)
		}
		res.Header().Set(echo.HeaderXRequestID, rid)
		return next(c)
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randString(8)
}

// randString generates a random string of length n
func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

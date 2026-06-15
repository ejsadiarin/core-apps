package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// NewLogger creates a new zerolog logger configured based on environment
func NewLogger(isProduction bool) zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if isProduction {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	return logger
}

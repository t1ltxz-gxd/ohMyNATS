package zerolog

import (
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func LoggerBuilder(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var logger *zerolog.Logger
	switch strings.ToLower(env) {
	case envLocal:
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			Level(zerolog.TraceLevel).With().Timestamp().Logger()
		return logger
	case envDev:
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			Level(zerolog.DebugLevel).With().Timestamp().Logger()
		return logger
	case envProd:
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			Level(zerolog.InfoLevel).With().Timestamp().Logger()
		return logger
	}
	return *logger
}

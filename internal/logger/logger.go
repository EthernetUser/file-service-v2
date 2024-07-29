package logger

import (
	"log/slog"
	"os"
)

const (
	EnvironmentProd       = "production"
	EnvironmentDevelopment = "development"
	EnvironmentLocal      = "local"
)

func NewLogger(environment string) *slog.Logger {
	var logger *slog.Logger

	switch environment {
	case EnvironmentProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	case EnvironmentDevelopment:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case EnvironmentLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return logger
}

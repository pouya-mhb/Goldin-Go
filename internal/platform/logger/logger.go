package logger

import (
	"log/slog"
	"os"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
)

// New creates the application's logger.
func New(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: parseLevel(cfg.Log.Level),
	}

	var handler slog.Handler

	switch cfg.App.Environment {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, opts)

	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler).With(
		slog.String("service", cfg.App.Name),
		slog.String("environment", cfg.App.Environment),
		slog.String("version", cfg.App.Version),
	)
}

func parseLevel(level string) slog.Level {
	switch level {

	case "debug":
		return slog.LevelDebug

	case "warn":
		return slog.LevelWarn

	case "error":
		return slog.LevelError

	default:
		return slog.LevelInfo
	}
}

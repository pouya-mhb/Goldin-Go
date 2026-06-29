package bootstrap

import (
	"database/sql"
	"log/slog"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
)

// App contains the application's infrastructure dependencies.
//
// It is created only during application startup and should never be used
// as a service locator throughout the codebase.
type App struct {
	Config *config.Config

	Infra *Infrastructure
}

// Infrastructure contains shared infrastructure services.
type Infrastructure struct {
	Logger *slog.Logger
	DB     *sql.DB

	// Redis  *redis.Client
	// Kafka  *kafka.Writer
	// Tracer trace.TracerProvider
}

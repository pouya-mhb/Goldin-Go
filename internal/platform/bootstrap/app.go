package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
	platformhttp "github.com/pouya-mhb/Goldin-Go/internal/platform/http"
)

// App contains the application's infrastructure dependencies.
//
// It is created only during application startup and should never be used
// as a service locator throughout the codebase.
type App struct {
	Config *config.Config

	Infra   *Infrastructure
	Modules *Modules
}

// Infrastructure contains shared infrastructure services.
type Infrastructure struct {
	Logger *slog.Logger
	DB     *sql.DB
	HTTP   *platformhttp.Server

	// Redis  *redis.Client
	// Kafka  *kafka.Writer
	// Tracer trace.TracerProvider
}

// Modules contains bounded context modules.
type Modules struct {
	Identity *identity.Module
}

// Run starts the application and blocks until the context is canceled.
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		errCh <- a.Infra.HTTP.Start()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.Infra.HTTP.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown application: %w", err)
		}

		if a.Infra.DB != nil {
			if err := a.Infra.DB.Close(); err != nil {
				return fmt.Errorf("close database: %w", err)
			}
		}

		a.Infra.Logger.Info("Goldin API stopped")

		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("run application: %w", err)
		}

		return nil
	}
}

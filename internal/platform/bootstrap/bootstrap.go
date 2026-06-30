package bootstrap

import (
	"fmt"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
	platformhttp "github.com/pouya-mhb/Goldin-Go/internal/platform/http"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/logger"
)

// Build constructs the application's dependency graph.
func Build() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	log := logger.New(cfg)
	router := platformhttp.NewRouter(log)
	httpServer := platformhttp.New(cfg.Server, log, router)

	app := &App{
		Config: cfg,
		Infra: &Infrastructure{
			Logger: log,
			HTTP:   httpServer,
		},
	}

	return app, nil
}

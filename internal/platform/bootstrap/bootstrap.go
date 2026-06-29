package bootstrap

import (
	"fmt"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/logger"
)

// Build constructs the application's dependency graph.
func Build() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration: %w", err)
	}

	log := logger.New(cfg)

	app := &App{
		Config: cfg,
		Infra: &Infrastructure{
			Logger: log,
		},
	}

	return app, nil
}

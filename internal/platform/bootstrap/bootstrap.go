package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity"
	identityhttp "github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/http"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/database"
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

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.OpenMySQL(dbCtx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	identityModule, err := identity.NewModule(db)
	if err != nil {
		return nil, fmt.Errorf("build identity module: %w", err)
	}

	router := platformhttp.NewRouter(
		log,
		identityhttp.WithRoutes(identityModule.RegisterUser),
	)
	httpServer := platformhttp.New(cfg.Server, log, router)

	app := &App{
		Config: cfg,
		Infra: &Infrastructure{
			Logger: log,
			DB:     db,
			HTTP:   httpServer,
		},
		Modules: &Modules{
			Identity: identityModule,
		},
	}

	return app, nil
}

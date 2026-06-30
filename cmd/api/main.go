package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := bootstrap.Build()
	if err != nil {
		slog.Error("build application", slog.Any("error", err))
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		app.Infra.Logger.Error("run application", slog.Any("error", err))
		os.Exit(1)
	}
}

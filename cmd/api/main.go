package main

import (
	"log"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/bootstrap"
)

func main() {
	app, err := bootstrap.Build()
	if err != nil {
		log.Fatal(err)
	}

	app.Infra.Logger.Info("Goldin API started")
}

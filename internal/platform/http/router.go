package http

import (
	nethttp "net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"
)

func NewRouter(logger *slog.Logger) nethttp.Handler {
	r := chi.NewRouter()

	registerMiddleware(r, logger)

	registerRoutes(r)

	return r
}

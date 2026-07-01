package http

import (
	nethttp "net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"
)

// RouterOption customizes platform router construction.
type RouterOption func(chi.Router)

// NewRouter constructs the application HTTP router.
func NewRouter(logger *slog.Logger, options ...RouterOption) nethttp.Handler {
	r := chi.NewRouter()

	registerMiddleware(r, logger)

	registerRoutes(r)

	for _, option := range options {
		option(r)
	}

	return r
}

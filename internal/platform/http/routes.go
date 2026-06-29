package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r chi.Router) {

	r.Get("/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)

		_, _ = w.Write([]byte("OK"))
	})

}

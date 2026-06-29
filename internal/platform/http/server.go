package http

import (
	"context"
	"fmt"
	nethttp "net/http"
	"time"

	"log/slog"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
)

type Server struct {
	http *nethttp.Server
}

func New(cfg config.ServerConfig, logger *slog.Logger, handler nethttp.Handler) *Server {
	server := &nethttp.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return &Server{
		http: server,
	}
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

package http

import (
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"time"

	"log/slog"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
)

type Server struct {
	http *nethttp.Server
	log  *slog.Logger
}

// New constructs an HTTP server.
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
		log:  logger,
	}
}

// Addr returns the configured server listen address.
func (s *Server) Addr() string {
	return s.http.Addr
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	s.log.Info("starting HTTP server", slog.String("address", s.http.Addr))

	if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
		return fmt.Errorf("listen and serve HTTP: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	return nil
}

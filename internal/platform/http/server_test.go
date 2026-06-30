package http_test

import (
	"context"
	"log/slog"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
	platformhttp "github.com/pouya-mhb/Goldin-Go/internal/platform/http"
)

func TestServerAddr(t *testing.T) {
	t.Parallel()

	server := platformhttp.New(
		config.ServerConfig{Host: "127.0.0.1", Port: 8080},
		slog.Default(),
		nethttp.NewServeMux(),
	)

	if server.Addr() != "127.0.0.1:8080" {
		t.Fatalf("expected address %q, got %q", "127.0.0.1:8080", server.Addr())
	}
}

func TestServerShutdownWithoutStart(t *testing.T) {
	t.Parallel()

	server := platformhttp.New(
		config.ServerConfig{Host: "127.0.0.1", Port: 8080},
		slog.Default(),
		nethttp.NewServeMux(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("expected shutdown without start to succeed, got %v", err)
	}
}

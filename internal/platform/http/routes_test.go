package http_test

import (
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	platformhttp "github.com/pouya-mhb/Goldin-Go/internal/platform/http"
)

func TestHealthRoute(t *testing.T) {
	t.Parallel()

	router := platformhttp.NewRouter(slog.Default())
	request := httptest.NewRequest(nethttp.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("expected status %d, got %d", nethttp.StatusOK, response.Code)
	}

	if response.Body.String() != "OK" {
		t.Fatalf("expected body %q, got %q", "OK", response.Body.String())
	}
}

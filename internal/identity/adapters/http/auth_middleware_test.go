package http_test

import (
	"context"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	identityhttp "github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/http"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/token"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestRequireAuthentication(t *testing.T) {
	t.Parallel()

	verifiedToken := token.VerifiedToken{
		UserID: valueobject.NewUserID(),
		Email:  mustEmail(t),
	}
	verifier := &fakeAccessTokenVerifier{verifiedToken: verifiedToken}
	called := false
	next := nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		called = true

		authenticatedUser, ok := identityhttp.AuthenticatedUserFromContext(r.Context())
		if !ok {
			t.Fatal("expected authenticated user in context")
		}

		if authenticatedUser.UserID != verifiedToken.UserID {
			t.Fatal("expected authenticated user id to match")
		}

		w.WriteHeader(nethttp.StatusNoContent)
	})

	handler := identityhttp.RequireAuthentication(verifier)(next)
	request := httptest.NewRequest(nethttp.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer access-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != nethttp.StatusNoContent {
		t.Fatalf("expected status %d, got %d", nethttp.StatusNoContent, response.Code)
	}

	if !called {
		t.Fatal("expected next handler to be called")
	}

	if verifier.tokenValue != "access-token" {
		t.Fatalf("expected token value %q, got %q", "access-token", verifier.tokenValue)
	}
}

func TestRequireAuthenticationFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header string
		err    error
	}{
		{
			name: "missing authorization",
		},
		{
			name:   "wrong scheme",
			header: "Basic abc",
		},
		{
			name:   "empty bearer token",
			header: "Bearer ",
		},
		{
			name:   "verifier rejects token",
			header: "Bearer invalid-token",
			err:    token.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			verifier := &fakeAccessTokenVerifier{err: tt.err}
			next := nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
				t.Fatal("next handler must not be called")
			})
			handler := identityhttp.RequireAuthentication(verifier)(next)
			request := httptest.NewRequest(nethttp.MethodGet, "/protected", nil)
			if tt.header != "" {
				request.Header.Set("Authorization", tt.header)
			}
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			if response.Code != nethttp.StatusUnauthorized {
				t.Fatalf("expected status %d, got %d", nethttp.StatusUnauthorized, response.Code)
			}
		})
	}
}

type fakeAccessTokenVerifier struct {
	verifiedToken token.VerifiedToken
	err           error
	tokenValue    string
}

func (v *fakeAccessTokenVerifier) VerifyAccessToken(ctx context.Context, tokenValue string) (token.VerifiedToken, error) {
	if err := ctx.Err(); err != nil {
		return token.VerifiedToken{}, err
	}

	v.tokenValue = tokenValue

	return v.verifiedToken, v.err
}

func mustEmail(t *testing.T) valueobject.Email {
	t.Helper()

	email, err := valueobject.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("create email: %v", err)
	}

	return email
}

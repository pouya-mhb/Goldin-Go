package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/token"
)

type authContextKey string

const authenticatedUserContextKey authContextKey = "authenticated_user"

// AccessTokenVerifier verifies access tokens for HTTP authentication.
type AccessTokenVerifier interface {
	VerifyAccessToken(ctx context.Context, tokenValue string) (token.VerifiedToken, error)
}

// RequireAuthentication requires a valid Bearer access token.
func RequireAuthentication(verifier AccessTokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenValue, err := bearerToken(r.Header.Get("Authorization"))
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "valid bearer token is required")
				return
			}

			verifiedToken, err := verifier.VerifyAccessToken(r.Context(), tokenValue)
			if err != nil {
				writeError(w, http.StatusUnauthorized, "unauthorized", "valid bearer token is required")
				return
			}

			ctx := context.WithValue(r.Context(), authenticatedUserContextKey, verifiedToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AuthenticatedUserFromContext returns the verified user from request context.
func AuthenticatedUserFromContext(ctx context.Context) (token.VerifiedToken, bool) {
	verifiedToken, ok := ctx.Value(authenticatedUserContextKey).(token.VerifiedToken)
	return verifiedToken, ok
}

func bearerToken(header string) (string, error) {
	const prefix = "Bearer "

	if !strings.HasPrefix(header, prefix) {
		return "", errors.New("authorization header must contain bearer token")
	}

	tokenValue := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if tokenValue == "" {
		return "", errors.New("bearer token is required")
	}

	return tokenValue, nil
}

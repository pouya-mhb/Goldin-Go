package token_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/token"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
)

func TestJWTIssuerIssueTokens(t *testing.T) {
	t.Parallel()

	cfg := config.JWTConfig{
		Secret:                      "local-development-secret-value-32chars",
		AccessTokenDurationMinutes:  15,
		RefreshTokenDurationMinutes: 43200,
	}
	issuer := token.NewJWTIssuer(cfg)
	now := time.Date(2026, time.July, 2, 12, 0, 0, 0, time.UTC)
	issuer.WithClock(func() time.Time { return now })

	userID := valueobject.NewUserID()
	email, err := valueobject.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("create email: %v", err)
	}

	tokens, err := issuer.IssueTokens(context.Background(), userID, email)
	if err != nil {
		t.Fatalf("issue tokens: %v", err)
	}

	if tokens.TokenType != "Bearer" {
		t.Fatalf("expected token type Bearer, got %q", tokens.TokenType)
	}

	if tokens.AccessTokenExpiresIn != 900 {
		t.Fatalf("expected access expiry seconds 900, got %d", tokens.AccessTokenExpiresIn)
	}

	if tokens.RefreshTokenExpiresIn != 2592000 {
		t.Fatalf("expected refresh expiry seconds 2592000, got %d", tokens.RefreshTokenExpiresIn)
	}

	assertTokenClaims(t, tokens.AccessToken, cfg.Secret, userID.String(), email.String(), "access")
	assertTokenClaims(t, tokens.RefreshToken, cfg.Secret, userID.String(), email.String(), "refresh")
}

func TestJWTIssuerIssueTokensHonorsCanceledContext(t *testing.T) {
	t.Parallel()

	issuer := token.NewJWTIssuer(config.JWTConfig{
		Secret:                      "local-development-secret-value-32chars",
		AccessTokenDurationMinutes:  15,
		RefreshTokenDurationMinutes: 43200,
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	tokens, err := issuer.IssueTokens(ctx, valueobject.NewUserID(), mustEmail(t))

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}

	if tokens.AccessToken != "" {
		t.Fatal("expected empty tokens")
	}
}

func assertTokenClaims(t *testing.T, tokenValue string, secret string, userID string, email string, tokenUse string) {
	t.Helper()

	parsedToken, err := jwt.Parse(tokenValue, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected map claims")
	}

	if claims["sub"] != userID {
		t.Fatalf("expected subject %q, got %v", userID, claims["sub"])
	}

	if claims["email"] != email {
		t.Fatalf("expected email %q, got %v", email, claims["email"])
	}

	if claims["token_use"] != tokenUse {
		t.Fatalf("expected token use %q, got %v", tokenUse, claims["token_use"])
	}
}

func mustEmail(t *testing.T) valueobject.Email {
	t.Helper()

	email, err := valueobject.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("create email: %v", err)
	}

	return email
}

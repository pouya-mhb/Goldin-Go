package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/ports/outbound"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
)

const bearerTokenType = "Bearer"

var (
	// ErrInvalidToken is returned when a token cannot be trusted.
	ErrInvalidToken = errors.New("invalid token")
)

// VerifiedToken contains identity claims extracted from a trusted token.
type VerifiedToken struct {
	UserID valueobject.UserID
	Email  valueobject.Email
}

// JWTIssuer issues signed JWT access and refresh tokens.
type JWTIssuer struct {
	secret          []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
	clock           func() time.Time
}

// VerifyAccessToken verifies an access token and extracts identity claims.
func (i *JWTIssuer) VerifyAccessToken(ctx context.Context, tokenValue string) (VerifiedToken, error) {
	if err := ctx.Err(); err != nil {
		return VerifiedToken{}, fmt.Errorf("verify token context: %w", err)
	}

	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}

		return i.secret, nil
	})
	if err != nil {
		return VerifiedToken{}, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if !parsedToken.Valid {
		return VerifiedToken{}, ErrInvalidToken
	}

	tokenUse, ok := claims["token_use"].(string)
	if !ok || tokenUse != "access" {
		return VerifiedToken{}, ErrInvalidToken
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return VerifiedToken{}, ErrInvalidToken
	}

	userID, err := valueobject.ParseUserID(subject)
	if err != nil {
		return VerifiedToken{}, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	emailClaim, ok := claims["email"].(string)
	if !ok {
		return VerifiedToken{}, ErrInvalidToken
	}

	email, err := valueobject.NewEmail(emailClaim)
	if err != nil {
		return VerifiedToken{}, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}

	if err := ctx.Err(); err != nil {
		return VerifiedToken{}, fmt.Errorf("verify token context: %w", err)
	}

	return VerifiedToken{
		UserID: userID,
		Email:  email,
	}, nil
}

// NewJWTIssuer constructs a JWTIssuer.
func NewJWTIssuer(cfg config.JWTConfig) *JWTIssuer {
	return &JWTIssuer{
		secret:          []byte(cfg.Secret),
		accessDuration:  time.Duration(cfg.AccessTokenDurationMinutes) * time.Minute,
		refreshDuration: time.Duration(cfg.RefreshTokenDurationMinutes) * time.Minute,
		clock:           time.Now,
	}
}

// IssueTokens issues access and refresh tokens for a user.
func (i *JWTIssuer) IssueTokens(ctx context.Context, userID valueobject.UserID, email valueobject.Email) (outbound.IssuedTokens, error) {
	if err := ctx.Err(); err != nil {
		return outbound.IssuedTokens{}, fmt.Errorf("issue tokens context: %w", err)
	}

	now := i.clock().UTC()

	accessToken, err := i.issueToken(userID, email, "access", now, i.accessDuration)
	if err != nil {
		return outbound.IssuedTokens{}, fmt.Errorf("issue access token: %w", err)
	}

	refreshToken, err := i.issueToken(userID, email, "refresh", now, i.refreshDuration)
	if err != nil {
		return outbound.IssuedTokens{}, fmt.Errorf("issue refresh token: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return outbound.IssuedTokens{}, fmt.Errorf("issue tokens context: %w", err)
	}

	return outbound.IssuedTokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		TokenType:             bearerTokenType,
		AccessTokenExpiresIn:  int64(i.accessDuration.Seconds()),
		RefreshTokenExpiresIn: int64(i.refreshDuration.Seconds()),
	}, nil
}

// WithClock replaces the token issuer clock for deterministic tests.
func (i *JWTIssuer) WithClock(clock func() time.Time) {
	i.clock = clock
}

func (i *JWTIssuer) issueToken(
	userID valueobject.UserID,
	email valueobject.Email,
	tokenUse string,
	issuedAt time.Time,
	duration time.Duration,
) (string, error) {
	claims := jwt.MapClaims{
		"sub":       userID.String(),
		"email":     email.String(),
		"token_use": tokenUse,
		"iat":       issuedAt.Unix(),
		"exp":       issuedAt.Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(i.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	return signedToken, nil
}

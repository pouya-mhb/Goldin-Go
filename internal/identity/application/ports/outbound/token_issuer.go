package outbound

import (
	"context"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

// TokenIssuer issues authentication tokens for an authenticated user.
type TokenIssuer interface {
	IssueTokens(ctx context.Context, userID valueobject.UserID, email valueobject.Email) (IssuedTokens, error)
}

// IssuedTokens contains access and refresh tokens.
type IssuedTokens struct {
	AccessToken           string
	RefreshToken          string
	TokenType             string
	AccessTokenExpiresIn  int64
	RefreshTokenExpiresIn int64
}

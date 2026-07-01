package outbound

import (
	"context"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

// PasswordVerifier verifies plaintext passwords against stored password hashes.
type PasswordVerifier interface {
	VerifyPassword(ctx context.Context, plaintext string, hash valueobject.PasswordHash) error
}

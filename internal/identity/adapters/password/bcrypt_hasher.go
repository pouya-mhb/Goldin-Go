package password

import (
	"context"
	"errors"
	"fmt"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrPasswordMismatch is returned when a password does not match a hash.
	ErrPasswordMismatch = errors.New("password mismatch")
)

// BcryptHasher hashes passwords with bcrypt.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher constructs a bcrypt password hasher.
func NewBcryptHasher(cost int) (*BcryptHasher, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, fmt.Errorf("bcrypt cost must be between %d and %d", bcrypt.MinCost, bcrypt.MaxCost)
	}

	return &BcryptHasher{cost: cost}, nil
}

// HashPassword hashes a plaintext password.
func (h *BcryptHasher) HashPassword(ctx context.Context, plaintext string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("hash password context: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), h.cost)
	if err != nil {
		return "", fmt.Errorf("generate bcrypt password hash: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("hash password context: %w", err)
	}

	return string(hash), nil
}

// VerifyPassword verifies a plaintext password against a stored hash.
func (h *BcryptHasher) VerifyPassword(ctx context.Context, plaintext string, hash valueobject.PasswordHash) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("verify password context: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash.String()), []byte(plaintext)); err != nil {
		return ErrPasswordMismatch
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("verify password context: %w", err)
	}

	return nil
}

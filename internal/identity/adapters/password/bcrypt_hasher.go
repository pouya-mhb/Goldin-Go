package password

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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

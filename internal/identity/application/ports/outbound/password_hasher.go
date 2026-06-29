package outbound

import "context"

// PasswordHasher hashes plaintext passwords before they enter the domain.
type PasswordHasher interface {
	HashPassword(ctx context.Context, plaintext string) (string, error)
}

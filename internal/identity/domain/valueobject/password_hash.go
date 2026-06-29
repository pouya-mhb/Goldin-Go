package valueobject

import "strings"

const maxPasswordHashLength = 255

// PasswordHash represents a hashed user password.
//
// The domain stores only the hash. Hashing and verification algorithms belong
// in application ports and adapters because they are technical policies.
type PasswordHash struct {
	value string
}

// NewPasswordHash validates a password hash before it enters the domain.
func NewPasswordHash(value string) (PasswordHash, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return PasswordHash{}, ErrEmptyPasswordHash
	}

	if len(trimmed) > maxPasswordHashLength {
		return PasswordHash{}, ErrPasswordHashTooLong
	}

	return PasswordHash{value: trimmed}, nil
}

// String returns the password hash value.
func (hash PasswordHash) String() string {
	return hash.value
}

// IsZero reports whether the password hash is unset.
func (hash PasswordHash) IsZero() bool {
	return hash.value == ""
}

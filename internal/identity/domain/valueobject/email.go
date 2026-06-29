package valueobject

import (
	"fmt"
	"net/mail"
	"strings"
)

const maxEmailLength = 254

// Email represents a normalized user email address.
type Email struct {
	value string
}

// NewEmail validates and normalizes an email address.
func NewEmail(value string) (Email, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return Email{}, ErrEmptyEmail
	}

	if len(normalized) > maxEmailLength {
		return Email{}, ErrEmailTooLong
	}

	address, err := mail.ParseAddress(normalized)
	if err != nil {
		return Email{}, fmt.Errorf("%w: %w", ErrInvalidEmail, err)
	}

	if address.Name != "" || address.Address != normalized {
		return Email{}, fmt.Errorf("%w: email must be a plain address", ErrInvalidEmail)
	}

	return Email{value: normalized}, nil
}

// String returns the normalized email address.
func (email Email) String() string {
	return email.value
}

// IsZero reports whether the email address is unset.
func (email Email) IsZero() bool {
	return email.value == ""
}

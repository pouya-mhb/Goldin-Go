package valueobject

import (
	"fmt"

	"github.com/google/uuid"
)

// UserID identifies a user across the Identity bounded context.
type UserID struct {
	value uuid.UUID
}

// NewUserID creates a new unique user identifier.
func NewUserID() UserID {
	return UserID{value: uuid.New()}
}

// ParseUserID validates and creates a UserID from its string representation.
func ParseUserID(value string) (UserID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return UserID{}, fmt.Errorf("parse user id: %w", err)
	}

	if id == uuid.Nil {
		return UserID{}, ErrEmptyUserID
	}

	return UserID{value: id}, nil
}

// String returns the canonical string representation of the user identifier.
func (id UserID) String() string {
	return id.value.String()
}

// IsZero reports whether the identifier is unset.
func (id UserID) IsZero() bool {
	return id.value == uuid.Nil
}

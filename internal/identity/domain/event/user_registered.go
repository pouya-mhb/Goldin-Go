package event

import (
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

// UserRegistered records that a user completed registration.
type UserRegistered struct {
	UserID     valueobject.UserID
	Email      valueobject.Email
	OccurredAt time.Time
}

// NewUserRegistered creates a UserRegistered domain event.
func NewUserRegistered(userID valueobject.UserID, email valueobject.Email, occurredAt time.Time) UserRegistered {
	return UserRegistered{
		UserID:     userID,
		Email:      email,
		OccurredAt: occurredAt.UTC(),
	}
}

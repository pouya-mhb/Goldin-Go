package domain

import (
	"errors"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/event"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

var (
	// ErrInvalidUserID is returned when a user aggregate has no identifier.
	ErrInvalidUserID = errors.New("user id is required")

	// ErrInvalidUserEmail is returned when a user aggregate has no email.
	ErrInvalidUserEmail = errors.New("user email is required")

	// ErrInvalidUserPasswordHash is returned when a user aggregate has no password hash.
	ErrInvalidUserPasswordHash = errors.New("user password hash is required")
)

// User is the aggregate root for identity accounts.
type User struct {
	id           valueobject.UserID
	email        valueobject.Email
	passwordHash valueobject.PasswordHash
	registeredAt time.Time
	events       []event.UserRegistered
}

// RegisterUser creates a new user aggregate and records a registration event.
func RegisterUser(
	id valueobject.UserID,
	email valueobject.Email,
	passwordHash valueobject.PasswordHash,
	registeredAt time.Time,
) (*User, error) {
	if id.IsZero() {
		return nil, ErrInvalidUserID
	}

	if email.IsZero() {
		return nil, ErrInvalidUserEmail
	}

	if passwordHash.IsZero() {
		return nil, ErrInvalidUserPasswordHash
	}

	if registeredAt.IsZero() {
		registeredAt = time.Now().UTC()
	}

	user := &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		registeredAt: registeredAt.UTC(),
	}

	user.record(event.NewUserRegistered(id, email, user.registeredAt))

	return user, nil
}

// RehydrateUser restores an existing user aggregate from persistence.
func RehydrateUser(
	id valueobject.UserID,
	email valueobject.Email,
	passwordHash valueobject.PasswordHash,
	registeredAt time.Time,
) (*User, error) {
	if id.IsZero() {
		return nil, ErrInvalidUserID
	}

	if email.IsZero() {
		return nil, ErrInvalidUserEmail
	}

	if passwordHash.IsZero() {
		return nil, ErrInvalidUserPasswordHash
	}

	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		registeredAt: registeredAt.UTC(),
	}, nil
}

// ID returns the user's identifier.
func (u *User) ID() valueobject.UserID {
	return u.id
}

// Email returns the user's email address.
func (u *User) Email() valueobject.Email {
	return u.email
}

// PasswordHash returns the user's password hash.
func (u *User) PasswordHash() valueobject.PasswordHash {
	return u.passwordHash
}

// RegisteredAt returns when the user registered.
func (u *User) RegisteredAt() time.Time {
	return u.registeredAt
}

// PullEvents returns pending domain events and clears them from the aggregate.
func (u *User) PullEvents() []event.UserRegistered {
	events := append([]event.UserRegistered(nil), u.events...)
	u.events = nil

	return events
}

func (u *User) record(event event.UserRegistered) {
	u.events = append(u.events, event)
}

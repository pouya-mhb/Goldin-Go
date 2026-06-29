package valueobject

import "errors"

var (
	// ErrEmailTooLong is returned when an email exceeds the maximum accepted length.
	ErrEmailTooLong = errors.New("email is too long")

	// ErrEmptyEmail is returned when an email is not provided.
	ErrEmptyEmail = errors.New("email is required")

	// ErrInvalidEmail is returned when an email is not syntactically valid.
	ErrInvalidEmail = errors.New("email is invalid")

	// ErrEmptyPasswordHash is returned when a password hash is not provided.
	ErrEmptyPasswordHash = errors.New("password hash is required")

	// ErrEmptyUserID is returned when a user identifier is empty.
	ErrEmptyUserID = errors.New("user id is required")

	// ErrPasswordHashTooLong is returned when a password hash exceeds the storage limit.
	ErrPasswordHashTooLong = errors.New("password hash is too long")
)

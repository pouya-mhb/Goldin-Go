package result

import "time"

// RegisterUser is returned after a user is registered.
type RegisterUser struct {
	UserID       string
	Email        string
	RegisteredAt time.Time
}

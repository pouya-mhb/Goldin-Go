package inbound

import (
	"context"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/result"
)

// RegisterUser handles user registration requests.
type RegisterUser interface {
	RegisterUser(ctx context.Context, cmd command.RegisterUser) (result.RegisterUser, error)
}

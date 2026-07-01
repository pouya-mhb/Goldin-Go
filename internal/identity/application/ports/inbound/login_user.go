package inbound

import (
	"context"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/result"
)

// LoginUser handles user login requests.
type LoginUser interface {
	LoginUser(ctx context.Context, cmd command.LoginUser) (result.LoginUser, error)
}

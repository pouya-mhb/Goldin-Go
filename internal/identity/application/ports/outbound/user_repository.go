package outbound

import (
	"context"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

// UserRepository persists and retrieves identity users.
type UserRepository interface {
	ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error)
	FindByEmail(ctx context.Context, email valueobject.Email) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
}

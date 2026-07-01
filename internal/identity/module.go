package identity

import (
	"database/sql"
	"fmt"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/password"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/persistence"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/ports/inbound"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/service"
	"golang.org/x/crypto/bcrypt"
)

const defaultBcryptCost = bcrypt.DefaultCost

// Module contains the Identity bounded context dependencies.
type Module struct {
	RegisterUser inbound.RegisterUser
}

// NewModule constructs the Identity module.
func NewModule(db *sql.DB) (*Module, error) {
	userRepository := persistence.NewMySQLUserRepository(db)

	passwordHasher, err := password.NewBcryptHasher(defaultBcryptCost)
	if err != nil {
		return nil, fmt.Errorf("create bcrypt password hasher: %w", err)
	}

	registration := service.NewRegistrationService(userRepository, passwordHasher)

	return &Module{
		RegisterUser: registration,
	}, nil
}

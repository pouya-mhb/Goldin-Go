package persistence

import (
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

type userRow struct {
	id           string
	email        string
	passwordHash string
	registeredAt time.Time
}

func (r userRow) toDomain() (*domain.User, error) {
	id, err := valueobject.ParseUserID(r.id)
	if err != nil {
		return nil, err
	}

	email, err := valueobject.NewEmail(r.email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := valueobject.NewPasswordHash(r.passwordHash)
	if err != nil {
		return nil, err
	}

	return domain.RehydrateUser(id, email, passwordHash, r.registeredAt)
}

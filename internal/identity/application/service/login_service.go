package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/persistence"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/ports/outbound"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/result"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

var (
	// ErrInvalidCredentials is returned when login credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// LoginService coordinates user login.
type LoginService struct {
	users     outbound.UserRepository
	passwords outbound.PasswordVerifier
}

// NewLoginService constructs a LoginService.
func NewLoginService(users outbound.UserRepository, passwords outbound.PasswordVerifier) *LoginService {
	return &LoginService{
		users:     users,
		passwords: passwords,
	}
}

// LoginUser authenticates a user with email and password credentials.
func (s *LoginService) LoginUser(ctx context.Context, cmd command.LoginUser) (result.LoginUser, error) {
	if err := ctx.Err(); err != nil {
		return result.LoginUser{}, fmt.Errorf("login context: %w", err)
	}

	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return result.LoginUser{}, ErrInvalidCredentials
	}

	user, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, persistence.ErrUserNotFound) {
		return result.LoginUser{}, ErrInvalidCredentials
	}

	if err != nil {
		return result.LoginUser{}, fmt.Errorf("find user by email: %w", err)
	}

	if err := s.passwords.VerifyPassword(ctx, cmd.Password, user.PasswordHash()); err != nil {
		return result.LoginUser{}, ErrInvalidCredentials
	}

	return result.LoginUser{
		UserID: user.ID().String(),
		Email:  user.Email().String(),
	}, nil
}

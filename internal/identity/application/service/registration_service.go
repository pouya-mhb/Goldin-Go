package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/ports/outbound"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/result"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

const minPasswordLength = 12

var (
	// ErrEmailAlreadyRegistered is returned when a user already owns the email.
	ErrEmailAlreadyRegistered = errors.New("email already registered")

	// ErrPasswordTooShort is returned when a password does not meet the minimum length.
	ErrPasswordTooShort = errors.New("password must be at least 12 characters")
)

// Clock returns the current time.
type Clock func() time.Time

// RegistrationService coordinates user registration.
type RegistrationService struct {
	users          outbound.UserRepository
	passwords      outbound.PasswordHasher
	clock          Clock
	newUserID      func() valueobject.UserID
	minimumPassLen int
}

// RegistrationOption customizes RegistrationService construction.
type RegistrationOption func(*RegistrationService)

// NewRegistrationService constructs a RegistrationService.
func NewRegistrationService(
	users outbound.UserRepository,
	passwords outbound.PasswordHasher,
	options ...RegistrationOption,
) *RegistrationService {
	service := &RegistrationService{
		users:          users,
		passwords:      passwords,
		clock:          time.Now,
		newUserID:      valueobject.NewUserID,
		minimumPassLen: minPasswordLength,
	}

	for _, option := range options {
		option(service)
	}

	return service
}

// RegisterUser registers a new identity user.
func (s *RegistrationService) RegisterUser(ctx context.Context, cmd command.RegisterUser) (result.RegisterUser, error) {
	if err := ctx.Err(); err != nil {
		return result.RegisterUser{}, fmt.Errorf("registration context: %w", err)
	}

	email, err := valueobject.NewEmail(cmd.Email)
	if err != nil {
		return result.RegisterUser{}, fmt.Errorf("create email: %w", err)
	}

	if len(strings.TrimSpace(cmd.Password)) < s.minimumPassLen {
		return result.RegisterUser{}, ErrPasswordTooShort
	}

	exists, err := s.users.ExistsByEmail(ctx, email)
	if err != nil {
		return result.RegisterUser{}, fmt.Errorf("check email uniqueness: %w", err)
	}

	if exists {
		return result.RegisterUser{}, ErrEmailAlreadyRegistered
	}

	hashedPassword, err := s.passwords.HashPassword(ctx, cmd.Password)
	if err != nil {
		return result.RegisterUser{}, fmt.Errorf("hash password: %w", err)
	}

	passwordHash, err := valueobject.NewPasswordHash(hashedPassword)
	if err != nil {
		return result.RegisterUser{}, fmt.Errorf("create password hash: %w", err)
	}

	user, err := domain.RegisterUser(s.newUserID(), email, passwordHash, s.clock())
	if err != nil {
		return result.RegisterUser{}, fmt.Errorf("register user domain: %w", err)
	}

	if err := s.users.Save(ctx, user); err != nil {
		return result.RegisterUser{}, fmt.Errorf("save user: %w", err)
	}

	return result.RegisterUser{
		UserID:       user.ID().String(),
		Email:        user.Email().String(),
		RegisteredAt: user.RegisteredAt(),
	}, nil
}

// WithClock replaces the service clock for deterministic tests or controlled runtimes.
func WithClock(clock Clock) RegistrationOption {
	return func(s *RegistrationService) {
		s.clock = clock
	}
}

// WithUserIDGenerator replaces the user identifier generator.
func WithUserIDGenerator(newUserID func() valueobject.UserID) RegistrationOption {
	return func(s *RegistrationService) {
		s.newUserID = newUserID
	}
}

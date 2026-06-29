package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/service"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestRegistrationServiceRegisterUser(t *testing.T) {
	t.Parallel()

	userID, err := valueobject.ParseUserID("21f3483a-176e-40e4-9d77-6d15fcb675d8")
	if err != nil {
		t.Fatalf("parse user id: %v", err)
	}

	registeredAt := time.Date(2026, time.June, 30, 10, 0, 0, 0, time.UTC)
	users := &fakeUserRepository{}
	passwords := &fakePasswordHasher{hash: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS"}
	registration := service.NewRegistrationService(
		users,
		passwords,
		service.WithClock(func() time.Time { return registeredAt }),
		service.WithUserIDGenerator(func() valueobject.UserID { return userID }),
	)

	result, err := registration.RegisterUser(context.Background(), command.RegisterUser{
		Email:    " USER@Example.COM ",
		Password: "correct horse battery staple",
	})

	if err != nil {
		t.Fatalf("register user: %v", err)
	}

	if result.UserID != userID.String() {
		t.Fatalf("expected user id %q, got %q", userID.String(), result.UserID)
	}

	if result.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", result.Email)
	}

	if !result.RegisteredAt.Equal(registeredAt) {
		t.Fatalf("expected registered at %s, got %s", registeredAt, result.RegisteredAt)
	}

	if users.saved == nil {
		t.Fatal("expected user to be saved")
	}

	if passwords.received != "correct horse battery staple" {
		t.Fatalf("expected hasher to receive plaintext password, got %q", passwords.received)
	}
}

func TestRegistrationServiceRegisterUserFailures(t *testing.T) {
	t.Parallel()

	repositoryErr := errors.New("repository unavailable")
	hasherErr := errors.New("hasher unavailable")

	tests := []struct {
		name      string
		command   command.RegisterUser
		users     *fakeUserRepository
		passwords *fakePasswordHasher
		wantErr   error
	}{
		{
			name: "rejects invalid email",
			command: command.RegisterUser{
				Email:    "invalid-email",
				Password: "correct horse battery staple",
			},
			users:     &fakeUserRepository{},
			passwords: &fakePasswordHasher{},
			wantErr:   valueobject.ErrInvalidEmail,
		},
		{
			name: "rejects short password",
			command: command.RegisterUser{
				Email:    "user@example.com",
				Password: "too-short",
			},
			users:     &fakeUserRepository{},
			passwords: &fakePasswordHasher{},
			wantErr:   service.ErrPasswordTooShort,
		},
		{
			name: "rejects duplicate email",
			command: command.RegisterUser{
				Email:    "user@example.com",
				Password: "correct horse battery staple",
			},
			users:     &fakeUserRepository{exists: true},
			passwords: &fakePasswordHasher{},
			wantErr:   service.ErrEmailAlreadyRegistered,
		},
		{
			name: "wraps repository uniqueness failure",
			command: command.RegisterUser{
				Email:    "user@example.com",
				Password: "correct horse battery staple",
			},
			users:     &fakeUserRepository{existsErr: repositoryErr},
			passwords: &fakePasswordHasher{},
			wantErr:   repositoryErr,
		},
		{
			name: "wraps password hashing failure",
			command: command.RegisterUser{
				Email:    "user@example.com",
				Password: "correct horse battery staple",
			},
			users:     &fakeUserRepository{},
			passwords: &fakePasswordHasher{err: hasherErr},
			wantErr:   hasherErr,
		},
		{
			name: "wraps invalid hash failure",
			command: command.RegisterUser{
				Email:    "user@example.com",
				Password: "correct horse battery staple",
			},
			users:     &fakeUserRepository{},
			passwords: &fakePasswordHasher{hash: ""},
			wantErr:   valueobject.ErrEmptyPasswordHash,
		},
		{
			name: "wraps save failure",
			command: command.RegisterUser{
				Email:    "user@example.com",
				Password: "correct horse battery staple",
			},
			users:     &fakeUserRepository{saveErr: repositoryErr},
			passwords: &fakePasswordHasher{hash: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS"},
			wantErr:   repositoryErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			registration := service.NewRegistrationService(tt.users, tt.passwords)

			_, err := registration.RegisterUser(context.Background(), tt.command)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestRegistrationServiceRegisterUserHonorsCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	registration := service.NewRegistrationService(&fakeUserRepository{}, &fakePasswordHasher{})

	_, err := registration.RegisterUser(ctx, command.RegisterUser{
		Email:    "user@example.com",
		Password: "correct horse battery staple",
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

type fakeUserRepository struct {
	exists    bool
	existsErr error
	saveErr   error
	saved     *domain.User
}

func (r *fakeUserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	return r.exists, r.existsErr
}

func (r *fakeUserRepository) Save(ctx context.Context, user *domain.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.saved = user

	return r.saveErr
}

type fakePasswordHasher struct {
	hash     string
	err      error
	received string
}

func (h *fakePasswordHasher) HashPassword(ctx context.Context, plaintext string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	h.received = plaintext

	return h.hash, h.err
}

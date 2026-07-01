package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/persistence"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/service"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestLoginServiceLoginUser(t *testing.T) {
	t.Parallel()

	user := mustUser(t)
	users := &fakeLoginUserRepository{user: user}
	passwords := &fakePasswordVerifier{}
	login := service.NewLoginService(users, passwords)

	result, err := login.LoginUser(context.Background(), command.LoginUser{
		Email:    " USER@Example.COM ",
		Password: "correct horse battery staple",
	})

	if err != nil {
		t.Fatalf("login user: %v", err)
	}

	if result.UserID != user.ID().String() {
		t.Fatalf("expected user id %q, got %q", user.ID().String(), result.UserID)
	}

	if result.Email != user.Email().String() {
		t.Fatalf("expected email %q, got %q", user.Email().String(), result.Email)
	}

	if users.email.String() != "user@example.com" {
		t.Fatalf("expected normalized lookup email, got %q", users.email.String())
	}

	if passwords.plaintext != "correct horse battery staple" {
		t.Fatal("expected plaintext password to be passed to verifier")
	}

	if passwords.hash != user.PasswordHash() {
		t.Fatal("expected stored password hash to be passed to verifier")
	}
}

func TestLoginServiceLoginUserFailures(t *testing.T) {
	t.Parallel()

	repositoryErr := errors.New("repository unavailable")
	verifierErr := errors.New("password mismatch")

	tests := []struct {
		name      string
		command   command.LoginUser
		users     *fakeLoginUserRepository
		passwords *fakePasswordVerifier
		wantErr   error
	}{
		{
			name: "invalid email",
			command: command.LoginUser{
				Email:    "invalid-email",
				Password: "correct horse battery staple",
			},
			users:     &fakeLoginUserRepository{},
			passwords: &fakePasswordVerifier{},
			wantErr:   service.ErrInvalidCredentials,
		},
		{
			name: "unknown email",
			command: command.LoginUser{
				Email:    "user@example.com",
				Password: "correct horse battery staple",
			},
			users:     &fakeLoginUserRepository{err: persistence.ErrUserNotFound},
			passwords: &fakePasswordVerifier{},
			wantErr:   service.ErrInvalidCredentials,
		},
		{
			name: "repository failure",
			command: command.LoginUser{
				Email:    "user@example.com",
				Password: "correct horse battery staple",
			},
			users:     &fakeLoginUserRepository{err: repositoryErr},
			passwords: &fakePasswordVerifier{},
			wantErr:   repositoryErr,
		},
		{
			name: "wrong password",
			command: command.LoginUser{
				Email:    "user@example.com",
				Password: "wrong password",
			},
			users:     &fakeLoginUserRepository{user: mustUser(t)},
			passwords: &fakePasswordVerifier{err: verifierErr},
			wantErr:   service.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			login := service.NewLoginService(tt.users, tt.passwords)

			_, err := login.LoginUser(context.Background(), tt.command)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestLoginServiceLoginUserHonorsCanceledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	login := service.NewLoginService(&fakeLoginUserRepository{}, &fakePasswordVerifier{})

	_, err := login.LoginUser(ctx, command.LoginUser{
		Email:    "user@example.com",
		Password: "correct horse battery staple",
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

type fakeLoginUserRepository struct {
	user  *domain.User
	err   error
	email valueobject.Email
}

func (r *fakeLoginUserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	return false, errors.New("not implemented by login test fake")
}

func (r *fakeLoginUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*domain.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.email = email

	return r.user, r.err
}

func (r *fakeLoginUserRepository) Save(ctx context.Context, user *domain.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	return errors.New("not implemented by login test fake")
}

type fakePasswordVerifier struct {
	err       error
	plaintext string
	hash      valueobject.PasswordHash
}

func (v *fakePasswordVerifier) VerifyPassword(ctx context.Context, plaintext string, hash valueobject.PasswordHash) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	v.plaintext = plaintext
	v.hash = hash

	return v.err
}

func mustUser(t *testing.T) *domain.User {
	t.Helper()

	id := valueobject.NewUserID()
	email, err := valueobject.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("create email: %v", err)
	}

	passwordHash, err := valueobject.NewPasswordHash("$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS")
	if err != nil {
		t.Fatalf("create password hash: %v", err)
	}

	user, err := domain.RehydrateUser(id, email, passwordHash, time.Now())
	if err != nil {
		t.Fatalf("rehydrate user: %v", err)
	}

	return user
}

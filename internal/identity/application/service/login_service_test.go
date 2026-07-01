package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/persistence"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/ports/outbound"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/service"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestLoginServiceLoginUser(t *testing.T) {
	t.Parallel()

	user := mustUser(t)
	users := &fakeLoginUserRepository{user: user}
	passwords := &fakePasswordVerifier{}
	tokens := &fakeTokenIssuer{
		accessToken:           "access-token",
		refreshToken:          "refresh-token",
		tokenType:             "Bearer",
		accessTokenExpiresIn:  900,
		refreshTokenExpiresIn: 2592000,
	}
	login := service.NewLoginService(users, passwords, tokens)

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

	if tokens.userID != user.ID() {
		t.Fatal("expected token issuer to receive user id")
	}

	if result.AccessToken != "access-token" {
		t.Fatalf("expected access token, got %q", result.AccessToken)
	}

	if result.RefreshToken != "refresh-token" {
		t.Fatalf("expected refresh token, got %q", result.RefreshToken)
	}
}

func TestLoginServiceLoginUserFailures(t *testing.T) {
	t.Parallel()

	repositoryErr := errors.New("repository unavailable")
	verifierErr := errors.New("password mismatch")
	tokenErr := errors.New("token issuer unavailable")

	tests := []struct {
		name      string
		command   command.LoginUser
		users     *fakeLoginUserRepository
		passwords *fakePasswordVerifier
		tokens    *fakeTokenIssuer
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
			tokens:    &fakeTokenIssuer{},
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
			tokens:    &fakeTokenIssuer{},
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
			tokens:    &fakeTokenIssuer{},
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
			tokens:    &fakeTokenIssuer{},
			wantErr:   service.ErrInvalidCredentials,
		},
		{
			name: "token issuer failure",
			command: command.LoginUser{
				Email:    "user@example.com",
				Password: "correct horse battery staple",
			},
			users:     &fakeLoginUserRepository{user: mustUser(t)},
			passwords: &fakePasswordVerifier{},
			tokens:    &fakeTokenIssuer{err: tokenErr},
			wantErr:   tokenErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			login := service.NewLoginService(tt.users, tt.passwords, tt.tokens)

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

	login := service.NewLoginService(&fakeLoginUserRepository{}, &fakePasswordVerifier{}, &fakeTokenIssuer{})

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

type fakeTokenIssuer struct {
	err                   error
	userID                valueobject.UserID
	email                 valueobject.Email
	accessToken           string
	refreshToken          string
	tokenType             string
	accessTokenExpiresIn  int64
	refreshTokenExpiresIn int64
}

func (i *fakeTokenIssuer) IssueTokens(ctx context.Context, userID valueobject.UserID, email valueobject.Email) (outbound.IssuedTokens, error) {
	if err := ctx.Err(); err != nil {
		return outbound.IssuedTokens{}, err
	}

	i.userID = userID
	i.email = email

	return outbound.IssuedTokens{
		AccessToken:           i.accessToken,
		RefreshToken:          i.refreshToken,
		TokenType:             i.tokenType,
		AccessTokenExpiresIn:  i.accessTokenExpiresIn,
		RefreshTokenExpiresIn: i.refreshTokenExpiresIn,
	}, i.err
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

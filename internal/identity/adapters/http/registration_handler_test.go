package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"

	identityhttp "github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/http"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/command"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/result"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/application/service"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestRegistrationHandlerRegister(t *testing.T) {
	t.Parallel()

	registeredAt := time.Date(2026, time.July, 2, 10, 0, 0, 0, time.UTC)
	useCase := &fakeRegisterUser{
		result: result.RegisterUser{
			UserID:       "21f3483a-176e-40e4-9d77-6d15fcb675d8",
			Email:        "user@example.com",
			RegisteredAt: registeredAt,
		},
	}
	handler := identityhttp.NewRegistrationHandler(useCase)

	request := httptest.NewRequest(nethttp.MethodPost, "/identity/register", bytes.NewBufferString(`{"email":"user@example.com","password":"correct horse battery staple"}`))
	response := httptest.NewRecorder()

	handler.Register(response, request)

	if response.Code != nethttp.StatusCreated {
		t.Fatalf("expected status %d, got %d", nethttp.StatusCreated, response.Code)
	}

	if useCase.command.Email != "user@example.com" {
		t.Fatalf("expected email to be passed to use case, got %q", useCase.command.Email)
	}

	if useCase.command.Password != "correct horse battery staple" {
		t.Fatal("expected password to be passed to use case")
	}

	var body map[string]string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["user_id"] != useCase.result.UserID {
		t.Fatalf("expected user id %q, got %q", useCase.result.UserID, body["user_id"])
	}

	if body["email"] != useCase.result.Email {
		t.Fatalf("expected email %q, got %q", useCase.result.Email, body["email"])
	}
}

func TestRegistrationHandlerRegisterFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "malformed json",
			body:       `{`,
			wantStatus: nethttp.StatusBadRequest,
			wantCode:   "invalid_request",
		},
		{
			name:       "unknown field",
			body:       `{"email":"user@example.com","password":"correct horse battery staple","role":"admin"}`,
			wantStatus: nethttp.StatusBadRequest,
			wantCode:   "invalid_request",
		},
		{
			name:       "invalid email",
			body:       `{"email":"invalid-email","password":"correct horse battery staple"}`,
			err:        valueobject.ErrInvalidEmail,
			wantStatus: nethttp.StatusBadRequest,
			wantCode:   "invalid_registration",
		},
		{
			name:       "short password",
			body:       `{"email":"user@example.com","password":"short"}`,
			err:        service.ErrPasswordTooShort,
			wantStatus: nethttp.StatusBadRequest,
			wantCode:   "invalid_registration",
		},
		{
			name:       "duplicate email",
			body:       `{"email":"user@example.com","password":"correct horse battery staple"}`,
			err:        service.ErrEmailAlreadyRegistered,
			wantStatus: nethttp.StatusConflict,
			wantCode:   "email_already_registered",
		},
		{
			name:       "unexpected failure",
			body:       `{"email":"user@example.com","password":"correct horse battery staple"}`,
			err:        errors.New("database unavailable"),
			wantStatus: nethttp.StatusInternalServerError,
			wantCode:   "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := identityhttp.NewRegistrationHandler(&fakeRegisterUser{err: tt.err})
			request := httptest.NewRequest(nethttp.MethodPost, "/identity/register", bytes.NewBufferString(tt.body))
			response := httptest.NewRecorder()

			handler.Register(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, response.Code)
			}

			var body map[string]string
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if body["error"] != tt.wantCode {
				t.Fatalf("expected error code %q, got %q", tt.wantCode, body["error"])
			}
		})
	}
}

func TestLoginHandlerLogin(t *testing.T) {
	t.Parallel()

	useCase := &fakeLoginUser{
		result: result.LoginUser{
			UserID:                "21f3483a-176e-40e4-9d77-6d15fcb675d8",
			Email:                 "user@example.com",
			AccessToken:           "access-token",
			RefreshToken:          "refresh-token",
			TokenType:             "Bearer",
			AccessTokenExpiresIn:  900,
			RefreshTokenExpiresIn: 2592000,
		},
	}
	handler := identityhttp.NewLoginHandler(useCase)

	request := httptest.NewRequest(nethttp.MethodPost, "/identity/login", bytes.NewBufferString(`{"email":"user@example.com","password":"correct horse battery staple"}`))
	response := httptest.NewRecorder()

	handler.Login(response, request)

	if response.Code != nethttp.StatusOK {
		t.Fatalf("expected status %d, got %d", nethttp.StatusOK, response.Code)
	}

	if useCase.command.Email != "user@example.com" {
		t.Fatalf("expected email to be passed to use case, got %q", useCase.command.Email)
	}

	var body map[string]any
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["access_token"] != useCase.result.AccessToken {
		t.Fatal("expected access token in response")
	}

	if body["refresh_token"] != useCase.result.RefreshToken {
		t.Fatal("expected refresh token in response")
	}

	if body["token_type"] != "Bearer" {
		t.Fatalf("expected token type Bearer, got %v", body["token_type"])
	}
}

func TestLoginHandlerLoginFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "malformed json",
			body:       `{`,
			wantStatus: nethttp.StatusBadRequest,
			wantCode:   "invalid_request",
		},
		{
			name:       "unknown field",
			body:       `{"email":"user@example.com","password":"correct horse battery staple","role":"admin"}`,
			wantStatus: nethttp.StatusBadRequest,
			wantCode:   "invalid_request",
		},
		{
			name:       "invalid credentials",
			body:       `{"email":"user@example.com","password":"wrong password"}`,
			err:        service.ErrInvalidCredentials,
			wantStatus: nethttp.StatusUnauthorized,
			wantCode:   "invalid_credentials",
		},
		{
			name:       "unexpected failure",
			body:       `{"email":"user@example.com","password":"correct horse battery staple"}`,
			err:        errors.New("token issuer unavailable"),
			wantStatus: nethttp.StatusInternalServerError,
			wantCode:   "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := identityhttp.NewLoginHandler(&fakeLoginUser{err: tt.err})
			request := httptest.NewRequest(nethttp.MethodPost, "/identity/login", bytes.NewBufferString(tt.body))
			response := httptest.NewRecorder()

			handler.Login(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, response.Code)
			}

			var body map[string]string
			if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
				t.Fatalf("decode response: %v", err)
			}

			if body["error"] != tt.wantCode {
				t.Fatalf("expected error code %q, got %q", tt.wantCode, body["error"])
			}
		})
	}
}

type fakeRegisterUser struct {
	command command.RegisterUser
	result  result.RegisterUser
	err     error
}

func (f *fakeRegisterUser) RegisterUser(ctx context.Context, cmd command.RegisterUser) (result.RegisterUser, error) {
	if err := ctx.Err(); err != nil {
		return result.RegisterUser{}, err
	}

	f.command = cmd

	return f.result, f.err
}

type fakeLoginUser struct {
	command command.LoginUser
	result  result.LoginUser
	err     error
}

func (f *fakeLoginUser) LoginUser(ctx context.Context, cmd command.LoginUser) (result.LoginUser, error) {
	if err := ctx.Err(); err != nil {
		return result.LoginUser{}, err
	}

	f.command = cmd

	return f.result, f.err
}

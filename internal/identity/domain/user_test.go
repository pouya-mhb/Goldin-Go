package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestRegisterUser(t *testing.T) {
	t.Parallel()

	id := valueobject.NewUserID()
	email, err := valueobject.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("create email: %v", err)
	}

	passwordHash, err := valueobject.NewPasswordHash("$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS")
	if err != nil {
		t.Fatalf("create password hash: %v", err)
	}

	registeredAt := time.Date(2026, time.June, 29, 12, 0, 0, 0, time.UTC)

	user, err := domain.RegisterUser(id, email, passwordHash, registeredAt)

	if err != nil {
		t.Fatalf("register user: %v", err)
	}

	if user.ID() != id {
		t.Fatalf("expected user id %q, got %q", id.String(), user.ID().String())
	}

	if user.Email() != email {
		t.Fatalf("expected email %q, got %q", email.String(), user.Email().String())
	}

	if user.PasswordHash() != passwordHash {
		t.Fatal("expected password hash to match")
	}

	if !user.RegisteredAt().Equal(registeredAt) {
		t.Fatalf("expected registered at %s, got %s", registeredAt, user.RegisteredAt())
	}

	events := user.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].UserID != id {
		t.Fatalf("expected event user id %q, got %q", id.String(), events[0].UserID.String())
	}

	if events[0].Email != email {
		t.Fatalf("expected event email %q, got %q", email.String(), events[0].Email.String())
	}

	if !events[0].OccurredAt.Equal(registeredAt) {
		t.Fatalf("expected event time %s, got %s", registeredAt, events[0].OccurredAt)
	}

	if remaining := user.PullEvents(); len(remaining) != 0 {
		t.Fatalf("expected events to be cleared, got %d", len(remaining))
	}
}

func TestRegisterUserValidation(t *testing.T) {
	t.Parallel()

	id := valueobject.NewUserID()
	email, err := valueobject.NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("create email: %v", err)
	}

	passwordHash, err := valueobject.NewPasswordHash("$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS")
	if err != nil {
		t.Fatalf("create password hash: %v", err)
	}

	tests := []struct {
		name         string
		id           valueobject.UserID
		email        valueobject.Email
		passwordHash valueobject.PasswordHash
		wantErr      error
	}{
		{
			name:         "requires user id",
			email:        email,
			passwordHash: passwordHash,
			wantErr:      domain.ErrInvalidUserID,
		},
		{
			name:         "requires email",
			id:           id,
			passwordHash: passwordHash,
			wantErr:      domain.ErrInvalidUserEmail,
		},
		{
			name:    "requires password hash",
			id:      id,
			email:   email,
			wantErr: domain.ErrInvalidUserPasswordHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user, err := domain.RegisterUser(tt.id, tt.email, tt.passwordHash, time.Now())

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if user != nil {
				t.Fatal("expected user to be nil")
			}
		})
	}
}

func TestRehydrateUserDoesNotRecordEvents(t *testing.T) {
	t.Parallel()

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

	if events := user.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no events, got %d", len(events))
	}
}

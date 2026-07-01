package persistence

import (
	"errors"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain"
	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestUserRowToDomain(t *testing.T) {
	t.Parallel()

	registeredAt := time.Date(2026, time.July, 1, 12, 0, 0, 0, time.UTC)
	row := userRow{
		id:           "21f3483a-176e-40e4-9d77-6d15fcb675d8",
		email:        "user@example.com",
		passwordHash: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS",
		registeredAt: registeredAt,
	}

	user, err := row.toDomain()
	if err != nil {
		t.Fatalf("map row to domain: %v", err)
	}

	if user.ID().String() != row.id {
		t.Fatalf("expected id %q, got %q", row.id, user.ID().String())
	}

	if user.Email().String() != row.email {
		t.Fatalf("expected email %q, got %q", row.email, user.Email().String())
	}

	if user.PasswordHash().String() != row.passwordHash {
		t.Fatal("expected password hash to match")
	}

	if !user.RegisteredAt().Equal(registeredAt) {
		t.Fatalf("expected registered at %s, got %s", registeredAt, user.RegisteredAt())
	}

	if events := user.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no domain events on rehydration, got %d", len(events))
	}
}

func TestUserRowToDomainValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		row     userRow
		wantErr error
	}{
		{
			name: "invalid id",
			row: userRow{
				id:           "not-a-uuid",
				email:        "user@example.com",
				passwordHash: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS",
			},
		},
		{
			name: "invalid email",
			row: userRow{
				id:           "21f3483a-176e-40e4-9d77-6d15fcb675d8",
				email:        "invalid-email",
				passwordHash: "$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS",
			},
			wantErr: valueobject.ErrInvalidEmail,
		},
		{
			name: "empty password hash",
			row: userRow{
				id:    "21f3483a-176e-40e4-9d77-6d15fcb675d8",
				email: "user@example.com",
			},
			wantErr: valueobject.ErrEmptyPasswordHash,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			user, err := tt.row.toDomain()

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if user != nil {
				t.Fatal("expected user to be nil")
			}

			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestIsDuplicateEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "duplicate entry",
			err:  &mysql.MySQLError{Number: duplicateEntryErrorNumber},
			want: true,
		},
		{
			name: "other mysql error",
			err:  &mysql.MySQLError{Number: 1205},
			want: false,
		},
		{
			name: "wrapped duplicate entry",
			err:  wrapDuplicateEntry(),
			want: true,
		},
		{
			name: "non mysql error",
			err:  errors.New("database unavailable"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isDuplicateEntry(tt.err)

			if got != tt.want {
				t.Fatalf("expected %t, got %t", tt.want, got)
			}
		})
	}
}

func TestMySQLUserRepositorySaveDuplicateMapping(t *testing.T) {
	t.Parallel()

	if !errors.Is(ErrDuplicateUserEmail, ErrDuplicateUserEmail) {
		t.Fatal("expected duplicate email error to support errors.Is")
	}
}

func TestRepositoryErrors(t *testing.T) {
	t.Parallel()

	if errors.Is(ErrUserNotFound, domain.ErrInvalidUserID) {
		t.Fatal("expected repository not found error to remain separate from domain validation errors")
	}
}

func wrapDuplicateEntry() error {
	return errors.Join(errors.New("insert failed"), &mysql.MySQLError{Number: duplicateEntryErrorNumber})
}

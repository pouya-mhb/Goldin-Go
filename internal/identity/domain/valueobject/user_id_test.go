package valueobject_test

import (
	"testing"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestNewUserID(t *testing.T) {
	t.Parallel()

	id := valueobject.NewUserID()

	if id.IsZero() {
		t.Fatal("expected generated user id to be non-zero")
	}

	if id.String() == "" {
		t.Fatal("expected generated user id string to be non-empty")
	}
}

func TestParseUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid user id",
			value:   "21f3483a-176e-40e4-9d77-6d15fcb675d8",
			wantErr: false,
		},
		{
			name:    "invalid user id",
			value:   "not-a-uuid",
			wantErr: true,
		},
		{
			name:    "nil user id",
			value:   "00000000-0000-0000-0000-000000000000",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := valueobject.ParseUserID(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !id.IsZero() {
					t.Fatalf("expected zero user id, got %q", id.String())
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if id.String() != tt.value {
				t.Fatalf("expected %q, got %q", tt.value, id.String())
			}
		})
	}
}

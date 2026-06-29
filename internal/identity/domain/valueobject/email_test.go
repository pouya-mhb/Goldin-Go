package valueobject_test

import (
	"strings"
	"testing"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestNewEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		want      string
		wantError bool
	}{
		{
			name:  "normalizes email",
			value: " USER@Example.COM ",
			want:  "user@example.com",
		},
		{
			name:      "rejects empty email",
			value:     " ",
			wantError: true,
		},
		{
			name:      "rejects malformed email",
			value:     "invalid-email",
			wantError: true,
		},
		{
			name:      "rejects display name",
			value:     "User <user@example.com>",
			wantError: true,
		},
		{
			name:      "rejects overly long email",
			value:     strings.Repeat("a", 245) + "@example.com",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			email, err := valueobject.NewEmail(tt.value)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !email.IsZero() {
					t.Fatalf("expected zero email, got %q", email.String())
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if email.String() != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, email.String())
			}
		})
	}
}

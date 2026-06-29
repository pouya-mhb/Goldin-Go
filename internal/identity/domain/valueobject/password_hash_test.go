package valueobject_test

import (
	"strings"
	"testing"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/domain/valueobject"
)

func TestNewPasswordHash(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		want      string
		wantError bool
	}{
		{
			name:  "accepts hash",
			value: " $2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS ",
			want:  "$2a$10$7EqJtq98hPqEX7fNZaFWoOHiXfQ9NfxfXqg8zUspK90W7TDhMoyaS",
		},
		{
			name:      "rejects empty hash",
			value:     "",
			wantError: true,
		},
		{
			name:      "rejects overly long hash",
			value:     strings.Repeat("x", 256),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash, err := valueobject.NewPasswordHash(tt.value)

			if tt.wantError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if !hash.IsZero() {
					t.Fatalf("expected zero password hash, got %q", hash.String())
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if hash.String() != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, hash.String())
			}
		})
	}
}

package password_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pouya-mhb/Goldin-Go/internal/identity/adapters/password"
	"golang.org/x/crypto/bcrypt"
)

func TestNewBcryptHasher(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cost    int
		wantErr bool
	}{
		{
			name: "accepts minimum cost",
			cost: bcrypt.MinCost,
		},
		{
			name:    "rejects cost below minimum",
			cost:    bcrypt.MinCost - 1,
			wantErr: true,
		},
		{
			name:    "rejects cost above maximum",
			cost:    bcrypt.MaxCost + 1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hasher, err := password.NewBcryptHasher(tt.cost)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				if hasher != nil {
					t.Fatal("expected hasher to be nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if hasher == nil {
				t.Fatal("expected hasher")
			}
		})
	}
}

func TestBcryptHasherHashPassword(t *testing.T) {
	t.Parallel()

	hasher, err := password.NewBcryptHasher(bcrypt.MinCost)
	if err != nil {
		t.Fatalf("create hasher: %v", err)
	}

	hash, err := hasher.HashPassword(context.Background(), "correct horse battery staple")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	if hash == "correct horse battery staple" {
		t.Fatal("expected password hash to differ from plaintext")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("correct horse battery staple")); err != nil {
		t.Fatalf("expected hash to verify password: %v", err)
	}
}

func TestBcryptHasherHashPasswordHonorsCanceledContext(t *testing.T) {
	t.Parallel()

	hasher, err := password.NewBcryptHasher(bcrypt.MinCost)
	if err != nil {
		t.Fatalf("create hasher: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	hash, err := hasher.HashPassword(ctx, "correct horse battery staple")

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}

	if hash != "" {
		t.Fatalf("expected empty hash, got %q", hash)
	}
}

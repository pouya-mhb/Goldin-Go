package identity_test

import (
	"testing"

	"github.com/pouya-mhb/Goldin-Go/internal/identity"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
)

func TestNewModule(t *testing.T) {
	t.Parallel()

	module, err := identity.NewModule(nil, config.JWTConfig{
		Secret:                      "local-development-secret-value-32chars",
		AccessTokenDurationMinutes:  15,
		RefreshTokenDurationMinutes: 43200,
	})
	if err != nil {
		t.Fatalf("create identity module: %v", err)
	}

	if module.RegisterUser == nil {
		t.Fatal("expected register user use case to be wired")
	}

	if module.LoginUser == nil {
		t.Fatal("expected login user use case to be wired")
	}
}

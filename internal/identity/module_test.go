package identity_test

import (
	"testing"

	"github.com/pouya-mhb/Goldin-Go/internal/identity"
)

func TestNewModule(t *testing.T) {
	t.Parallel()

	module, err := identity.NewModule(nil)
	if err != nil {
		t.Fatalf("create identity module: %v", err)
	}

	if module.RegisterUser == nil {
		t.Fatal("expected register user use case to be wired")
	}
}

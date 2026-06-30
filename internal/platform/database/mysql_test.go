package database_test

import (
	"strings"
	"testing"

	"github.com/pouya-mhb/Goldin-Go/internal/platform/config"
	"github.com/pouya-mhb/Goldin-Go/internal/platform/database"
)

func TestMySQLDSN(t *testing.T) {
	t.Parallel()

	cfg := config.DatabaseConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "goldin",
		Password: "secret",
		Name:     "goldin",
	}

	dsn := database.MySQLDSN(cfg)

	expectedParts := []string{
		"goldin:secret@tcp(localhost:3306)/goldin",
		"charset=utf8mb4",
		"collation=utf8mb4_unicode_ci",
		"parseTime=true",
		"timeout=5s",
	}

	for _, part := range expectedParts {
		if !strings.Contains(dsn, part) {
			t.Fatalf("expected DSN %q to contain %q", dsn, part)
		}
	}
}

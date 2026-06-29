package config

import (
	"fmt"
)

var validLogLevels = map[string]struct{}{
	"debug": {},
	"info":  {},
	"warn":  {},
	"error": {},
}

// Validate validates the application configuration.
func Validate(cfg *Config) error {

	// Database

	if err := required("DB_HOST", cfg.Database.Host); err != nil {
		return err
	}

	if err := required("DB_USER", cfg.Database.User); err != nil {
		return err
	}

	if err := required("DB_NAME", cfg.Database.Name); err != nil {
		return err
	}

	// JWT

	if len(cfg.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	// Server

	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("SERVER_PORT must be between 1 and 65535")
	}

	// Database Port

	if cfg.Database.Port < 1 || cfg.Database.Port > 65535 {
		return fmt.Errorf("DB_PORT must be between 1 and 65535")
	}

	// Redis Port

	if cfg.Redis.Port < 1 || cfg.Redis.Port > 65535 {
		return fmt.Errorf("REDIS_PORT must be between 1 and 65535")
	}

	// Log Level

	if _, ok := validLogLevels[cfg.Log.Level]; !ok {
		return fmt.Errorf("invalid LOG_LEVEL: %s", cfg.Log.Level)
	}

	// Database Pool

	if cfg.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS must be greater than zero")
	}

	if cfg.Database.MaxIdleConns < 0 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS cannot be negative")
	}

	if cfg.Database.MaxIdleConns > cfg.Database.MaxOpenConns {
		return fmt.Errorf("DB_MAX_IDLE_CONNS cannot be greater than DB_MAX_OPEN_CONNS")
	}

	return nil
}

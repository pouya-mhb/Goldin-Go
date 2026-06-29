package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Load loads configuration from environment variables.
//
// In development it attempts to load a .env file.
// In production it relies entirely on environment variables.
func Load() (*Config, error) {
	// Ignore the error intentionally.
	// In production a .env file usually doesn't exist.
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "goldin"),
			Environment: getEnv("APP_ENV", "development"),
			Version:     getEnv("APP_VERSION", "0.0.1"),
		},

		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},

		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", ""),
			Port:         getEnvAsInt("DB_PORT", 3306),
			User:         getEnv("DB_USER", ""),
			Password:     getEnv("DB_PASSWORD", ""),
			Name:         getEnv("DB_NAME", ""),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", ""),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},

		JWT: JWTConfig{
			Secret:                      getEnv("JWT_SECRET", ""),
			AccessTokenDurationMinutes:  getEnvAsInt("JWT_ACCESS_DURATION", 15),
			RefreshTokenDurationMinutes: getEnvAsInt("JWT_REFRESH_DURATION", 43200),
		},

		Kafka: KafkaConfig{
			Brokers:  splitCSV(getEnv("KAFKA_BROKERS", "")),
			ClientID: getEnv("KAFKA_CLIENT_ID", "goldin-api"),
		},

		Log: LogConfig{
			Level: strings.ToLower(getEnv("LOG_LEVEL", "info")),
		},
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return i
}

func splitCSV(value string) []string {
	if value == "" {
		return []string{}
	}

	items := strings.Split(value, ",")

	var result []string

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}

	return result
}

func required(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", name)
	}

	return nil
}

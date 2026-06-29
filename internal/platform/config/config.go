package config

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Kafka    KafkaConfig
	Log      LogConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Version     string
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host string
	Port int

	User     string
	Password string
	Name     string

	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Host string
	Port int

	Password string
	DB       int
}

type JWTConfig struct {
	Secret string

	AccessTokenDurationMinutes  int
	RefreshTokenDurationMinutes int
}

type KafkaConfig struct {
	Brokers  []string
	ClientID string
}

type LogConfig struct {
	Level string
}

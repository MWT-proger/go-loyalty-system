package configs

import "os"

type AuthConfig struct {
	SecretKey string
}

type Config struct {
	HostServer  string `env:"RUN_ADDRESS"`
	LogLevel    string
	DatabaseDSN string `env:"DATABASE_DSN"`
	Auth        AuthConfig
}

var newConfig Config

// InitConfig() Присваивает локальной не импортируемой переменной newConfig базовые значения
// Вызывается один раз при старте проекта
func InitConfig() *Config {
	newConfig = Config{
		HostServer:  ":8000",
		LogLevel:    "info",
		DatabaseDSN: "",
		Auth:        AuthConfig{SecretKey: "supersecretkey"},
	}
	return &newConfig
}

// GetConfig() выводит не импортируемую переменную newConfig
func GetConfig() Config {
	return newConfig
}

// SetConfigFromEnv() Прсваевает полям значения из ENV
// Вызывается один раз при старте проекта
func SetConfigFromEnv() Config {

	if envBaseURLShortener := os.Getenv("RUN_ADDRESS"); envBaseURLShortener != "" {
		newConfig.HostServer = envBaseURLShortener
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		newConfig.LogLevel = envLogLevel
	}
	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		newConfig.DatabaseDSN = envDatabaseDSN
	}
	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		newConfig.Auth.SecretKey = envSecretKey
	}
	return newConfig
}

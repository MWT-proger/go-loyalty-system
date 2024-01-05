package configs

import (
	"flag"
	"os"
)

type AuthConfig struct {
	SecretKey string
}
type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	Debug            bool
}

type Config struct {
	HostServer           string `env:"RUN_ADDRESS"`
	LogLevel             string
	DatabaseDSN          string `env:"DATABASE_URI"`
	AccuralSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	Auth                 AuthConfig
	Cors                 CorsConfig
}

var newConfig Config

// initDefaultConfig() Присваивает локальной не импортируемой переменной newConfig базовые значения.
// Вызывается один раз при старте проекта.
func initDefaultConfig() {
	newConfig = Config{
		HostServer:           ":8000",
		LogLevel:             "info",
		DatabaseDSN:          "",
		AccuralSystemAddress: "localhost:7000",
		Auth:                 AuthConfig{SecretKey: "supersecretkey"},
		Cors: CorsConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"POST", "OPTIONS", "GET"},
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
				"Accept",
				"Origin",
				"Access-Control-Request-Method",
				"Access-Control-Request-Headers",
				"X-CSRF-Token",
			},
			AllowCredentials: true,
			Debug:            true,
		},
	}

}

// InitConfig() Инициализирует локальную не импортируемую переменную newConfig.
// Вызывает все доступные методы получения конфигов.
// Вызывается один раз при старте проекта
func InitConfig() *Config {
	initDefaultConfig()
	parseFlags()
	setConfigFromEnv()
	return &newConfig
}

// setConfigFromEnv() Прсваевает полям  глобальной переменной newConfig значения из ENV.
// Вызывается один раз при старте проекта.
func setConfigFromEnv() {

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
	if envAccuralSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccuralSystemAddress != "" {
		newConfig.AccuralSystemAddress = envAccuralSystemAddress
	}
}

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных глобальной переменной newConfig
func parseFlags() {

	flag.StringVar(&newConfig.HostServer, "a", newConfig.HostServer, "адрес и порт для запуска сервера")
	flag.StringVar(&newConfig.DatabaseDSN, "d", newConfig.DatabaseDSN, "строка с адресом подключения к БД")
	flag.StringVar(&newConfig.LogLevel, "l", "info", "уровень логирования")
	flag.StringVar(&newConfig.AccuralSystemAddress, "r", newConfig.AccuralSystemAddress, "адрес сервера системы начисления")
	flag.Parse()
}

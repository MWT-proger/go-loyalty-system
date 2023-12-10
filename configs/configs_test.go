package configs

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{name: "Тест 1", want: Config{
			HostServer:           ":1234",
			LogLevel:             "info",
			DatabaseDSN:          "",
			Auth:                 AuthConfig{SecretKey: "supersecretkey"},
			AccuralSystemAddress: "localhost:7000",
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
		}},
		{name: "Тест 2", want: Config{
			HostServer:           ":7777",
			LogLevel:             "debug",
			DatabaseDSN:          "",
			Auth:                 AuthConfig{SecretKey: "supersecretkey"},
			AccuralSystemAddress: "localhost:6000",
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
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newConfig = tt.want

			got := GetConfig()

			assert.Equal(t, got, tt.want, "GetConfig() не совпадает с ожидаемым")
		})
	}
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{name: "Тест 1", want: Config{
			HostServer:           ":8000",
			LogLevel:             "info",
			DatabaseDSN:          "",
			Auth:                 AuthConfig{SecretKey: "supersecretkey"},
			AccuralSystemAddress: "localhost:7000",
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
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitConfig()
			assert.Equal(t, newConfig, tt.want, "newConfig не совпадает с ожидаемым")
		})
	}
}

func TestGetConfigFromEnv(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{name: "Тест 1", want: Config{
			HostServer: ":7777",
			LogLevel:   "info",
			DatabaseDSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
				`localhost`, `postgres`, `postgres`, `testDB`),
			Auth:                 AuthConfig{SecretKey: "NewSuperSecretKeyTEEEEEEEEEEST"},
			AccuralSystemAddress: "example.ru",
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
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("RUN_ADDRESS", tt.want.HostServer)
			os.Setenv("LOG_LEVEL", tt.want.LogLevel)
			os.Setenv("ACCRUAL_SYSTEM_ADDRESS", tt.want.AccuralSystemAddress)
			os.Setenv("DATABASE_URI", tt.want.DatabaseDSN)
			os.Setenv("SECRET_KEY", tt.want.Auth.SecretKey)
			setConfigFromEnv()
			assert.Equal(t, newConfig, tt.want, "newConfig не совпадает с ожидаемым")
		})
	}
}

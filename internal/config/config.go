package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv    string
	Port      int
	DBHost    string
	DBPort    int
	DBUser    string
	DBPassword string
	DBName    string
	JWTSecret string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:    envOrDefault("APP_ENV", "development"),
		Port:      intOrDefault("PORT", 8080),
		DBHost:    envOrDefault("DB_HOST", "localhost"),
		DBPort:    intOrDefault("DB_PORT", 5432),
		DBUser:    envOrDefault("DB_USER", "chainhub"),
		DBPassword: envOrDefault("DB_PASSWORD", "chainhub"),
		DBName:    envOrDefault("DB_NAME", "chainhub"),
		JWTSecret: envOrDefault("JWT_SECRET", ""),
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func intOrDefault(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

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
	DBSSLMode string
	JWTSecret string
	MigrationsPath string
	RunMigrations  bool
	FrontendURL    string
}

func Load() (Config, error) {
	appEnv := envOrDefault("APP_ENV", "development")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	if dbSSLMode == "" {
		if appEnv == "production" {
			dbSSLMode = "require"
		} else {
			dbSSLMode = "disable"
		}
	}

	cfg := Config{
		AppEnv:    appEnv,
		Port:      intOrDefault("PORT", 8080),
		DBHost:    envOrDefault("DB_HOST", "localhost"),
		DBPort:    intOrDefault("DB_PORT", 5432),
		DBUser:    envOrDefault("DB_USER", "chainhub"),
		DBPassword: envOrDefault("DB_PASSWORD", "chainhub"),
		DBName:    envOrDefault("DB_NAME", "chainhub"),
		DBSSLMode: dbSSLMode,
		JWTSecret: envOrDefault("JWT_SECRET", ""),
		MigrationsPath: envOrDefault("MIGRATIONS_PATH", "file://migrations"),
		RunMigrations:  boolOrDefault("RUN_MIGRATIONS", true),
		FrontendURL:    envOrDefault("FRONTEND_URL", "http://localhost:3000"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

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

func boolOrDefault(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}

package db

import (
	"database/sql"
	"fmt"
	"time"

	"chainhub-api/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn(cfg))
	if err != nil {
		return nil, err
	}

	const (
		maxAttempts   = 10
		initialDelay  = 500 * time.Millisecond
		maxDelay      = 5 * time.Second
	)
	delay := initialDelay
	var pingErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := db.Ping(); err == nil {
			return db, nil
		} else {
			pingErr = err
		}
		if attempt < maxAttempts {
			time.Sleep(delay)
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
	}

	return nil, pingErr
}

func dsn(cfg config.Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)
}

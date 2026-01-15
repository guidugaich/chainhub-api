package handlers

import (
	"database/sql"

	"chainhub-api/internal/config"
)

type Handler struct {
	DB     *sql.DB
	Config config.Config
}

func New(db *sql.DB, cfg config.Config) *Handler {
	return &Handler{
		DB:     db,
		Config: cfg,
	}
}

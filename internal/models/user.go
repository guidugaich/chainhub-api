package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64
	Email        sql.NullString
	Username     sql.NullString
	PasswordHash sql.NullString
	WalletAddress sql.NullString
	CreatedAt    time.Time
}

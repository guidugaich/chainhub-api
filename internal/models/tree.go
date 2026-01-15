package models

import "time"

type Tree struct {
	ID        int64
	UserID    int64
	Title     string
	CreatedAt time.Time
}

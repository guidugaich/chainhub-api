package models

import "time"

type Link struct {
	ID        int64
	TreeID    int64
	Title     string
	URL       string
	Position  int
	IsActive  bool
	CreatedAt time.Time
}

package domain

import "time"

type Withdrawn struct {
	ID          int       `json:"-" db:"id"`
	Order       string    `json:"order" db:"order"`
	Sum         float64   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
	UserID      int       `json:"-" db:"user_id"`
}

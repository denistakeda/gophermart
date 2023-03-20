package domain

import "time"

type Withdrawn struct {
	ID          int       `db:"id"`
	OrderNumber string    `db:"order_number"`
	Sum         int       `db:"sum"`
	ProcessedAt time.Time `db:"processed_at"`
	UserID      int       `db:"user_id"`
}

type WithdrawnDisplay struct {
	OrderNumber string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (w Withdrawn) ToDisplay() WithdrawnDisplay {
	return WithdrawnDisplay{
		OrderNumber: w.OrderNumber,
		Sum:         float64(w.Sum) / 100,
		ProcessedAt: w.ProcessedAt,
	}
}

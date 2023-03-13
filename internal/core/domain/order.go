package domain

import "time"

type OrderStatus int

const (
	OrderStatus_New = iota
)

type Order struct {
	ID          int         `db:"id"`
	UserID      int         `db:"user_id"`
	OrderNumber int         `db:"order_number"`
	Status      OrderStatus `db:"status"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

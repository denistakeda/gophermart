package domain

import "time"

type OrderStatus string

const (
	OrderStatus_New        OrderStatus = "NEW"
	OrderStatus_Processing OrderStatus = "PROCESSING"
	OrderStatus_Invalid    OrderStatus = "INVALID"
	OrderStatus_Processed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID          int         `db:"id"`
	UserID      int         `db:"user_id"`
	OrderNumber int         `db:"order_number"`
	Status      OrderStatus `db:"status"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

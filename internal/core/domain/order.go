package domain

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID          int         `db:"id" json:"-"`
	UserID      int         `db:"user_id" json:"-"`
	OrderNumber int         `db:"order_number" json:"number"`
	Status      OrderStatus `db:"status" json:"status"`
	Accrual     int         `db:"accrual" json:"accrual,omitempty"`
	CreatedAt   time.Time   `db:"created_at" json:"uploaded_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"-"`
}

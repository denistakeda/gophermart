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
	Accrual     float64     `db:"accrual" json:"accrual,omitempty"`
	CreatedAt   time.Time   `db:"created_at" json:"uploaded_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"-"`
}

type UserBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

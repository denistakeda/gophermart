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
	OrderNumber string      `db:"order_number" json:"number"`
	Status      OrderStatus `db:"status" json:"status"`
	Accrual     int         `db:"accrual" json:"accrual,omitempty"`
	CreatedAt   time.Time   `db:"created_at" json:"uploaded_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"-"`
}

type OrderDisplay struct {
	OrderNumber string      `json:"number"`
	Status      OrderStatus `json:"status"`
	Accrual     float64     `json:"accrual,omitempty"`
	CreatedAt   time.Time   `json:"uploaded_at"`
}

func (o Order) ToDisplay() OrderDisplay {
	return OrderDisplay{
		OrderNumber: o.OrderNumber,
		Status:      o.Status,
		Accrual:     float64(o.Accrual) / 100,
		CreatedAt:   time.Time{},
	}
}

type UserBalance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

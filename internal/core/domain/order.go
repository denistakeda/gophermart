package domain

type OrderStatus int

const (
	OrderStatus_New = iota
)

type Order struct {
	ID          int
	UserID      int
	OrderNumber int
	Status      OrderStatus
}

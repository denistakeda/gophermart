package ports

import (
	"context"
	"gophermart/internal/core/domain"
)

type UserStore interface {
	AddNewUser(ctx context.Context, login, passwordHash string) error
	GetUser(ctx context.Context, login string) (domain.User, error)
}

type OrderStore interface {
	GetOrder(ctx context.Context, orderNumber int) (domain.Order, error)
	AddNewOrder(ctx context.Context, userID int, orderNumber int) error
	GetAllOrders(ctx context.Context, userID int) ([]domain.Order, error)
	GetAllNotFinished(ctx context.Context) ([]domain.Order, error)
	UpdateOrders(ctx context.Context, orders []domain.Order) error
}

type WithdrawnStore interface {
	GetAllWithdrawals(ctx context.Context, userID int) ([]domain.Withdrawn, error)
}

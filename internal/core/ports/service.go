package ports

import (
	"context"
	"gophermart/internal/core/domain"
)

type UserService interface {
	RegisterUser(ctx context.Context, login, password string) (string, error)
	LoginUser(ctx context.Context, login, password string) (string, error)
	AuthenticateUser(ctx context.Context, token string) (domain.User, error)
}

type OrderService interface {
	AddOrder(ctx context.Context, user *domain.User, orderNumber string) error
	GetAllOrders(ctx context.Context, user *domain.User) ([]domain.Order, error)
}

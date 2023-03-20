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
	GetUserBalance(ctx context.Context, user *domain.User) (domain.UserBalance, error)
	Withdraw(ctx context.Context, orderNumber string, sum float64, user *domain.User) error
	GetAllWithdrawals(ctx context.Context, user *domain.User) ([]domain.Withdrawn, error)
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

type AccrualService interface {
	CheckAccrual(orderNumber string) (AccrualResponse, error)
}

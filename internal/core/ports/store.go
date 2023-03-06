package ports

import (
	"context"
	"gophermart/internal/core/domain"
)

type UserStore interface {
	AddNewUser(ctx context.Context, login, passwordHash string) error
	GetUser(ctx context.Context, login string) (domain.User, error)
}

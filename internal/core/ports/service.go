package ports

import "context"

type UserService interface {
	RegisterUser(ctx context.Context, login, password string) (string, error)
}

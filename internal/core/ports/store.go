package ports

import "context"

type UserStore interface {
	AddNewUser(ctx context.Context, login, passwordHash string) error
}

package user_store

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

func (u *UserStore) AddNewUser(ctx context.Context, login, passwordHash string) error {
	if _, err := u.db.ExecContext(ctx, `
		insert into users(login, password, created_at)
		values ($1, $2, $3)
	`, login, passwordHash, time.Now()); err != nil {
		return errors.Wrap(err, "failed to insert user into a database")
	}

	return nil
}

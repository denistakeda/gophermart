package userstore

import (
	"context"
	"gophermart/internal/core/domain"
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

func (u *UserStore) GetUser(ctx context.Context, login string) (domain.User, error) {
	var user domain.User
	if err := u.db.GetContext(ctx, &user, `
		select login, password from users
		where login=$1
	`, login); err != nil {
		return user, errors.Wrapf(err, "failed to get user %s from the database", login)
	}

	return user, nil
}

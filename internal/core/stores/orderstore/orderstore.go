package orderstore

import (
	"context"
	"gophermart/internal/core/domain"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type OrderStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *OrderStore {
	return &OrderStore{db: db}
}

func (o *OrderStore) AddNewOrder(ctx context.Context, userID int, orderNumber int) error {
	if _, err := o.db.ExecContext(ctx, `
		insert into orders(user_id, order_number, status, created_at, updated_at)
		values ($1, $2, $3, $4, $5)
	`, userID, orderNumber, domain.OrderStatus_New, time.Now(), time.Now()); err != nil {
		return errors.Wrapf(err, "failed to insert order %d into a database", orderNumber)
	}

	return nil
}

func (o *OrderStore) GetOrder(ctx context.Context, orderNumber int) (domain.Order, error) {
	var order domain.Order
	if err := o.db.GetContext(ctx, &order, `
		select * from users
		where order_number=$1
	`, orderNumber); err != nil {
		return domain.Order{}, errors.Wrapf(err, "failed to get order %d from the database", orderNumber)
	}

	return order, nil
}

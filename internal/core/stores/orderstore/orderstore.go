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
	`, userID, orderNumber, domain.OrderStatusNew, time.Now(), time.Now()); err != nil {
		return errors.Wrapf(err, "failed to insert order %d into a database", orderNumber)
	}

	return nil
}

func (o *OrderStore) GetOrder(ctx context.Context, orderNumber int) (domain.Order, error) {
	var order domain.Order
	if err := o.db.GetContext(ctx, &order, `
		select * from orders
		where order_number=$1
	`, orderNumber); err != nil {
		return domain.Order{}, errors.Wrapf(err, "failed to get order %d from the database", orderNumber)
	}

	return order, nil
}

func (o *OrderStore) GetAllOrders(ctx context.Context, userID int) ([]domain.Order, error) {
	var orders []domain.Order
	if err := o.db.SelectContext(ctx, &orders, `
		select * from orders
		where user_id=$1
	`, userID); err != nil {
		return orders, errors.Wrapf(err, "failed to get list of orders for user with id %d", userID)
	}

	return orders, nil
}

func (o *OrderStore) GetAllNotFinished(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	if err := o.db.SelectContext(ctx, &orders, `
		select * from orders
		where status=$1 or status=$2
	`, domain.OrderStatusNew, domain.OrderStatusProcessing); err != nil {
		return orders, errors.Wrapf(err, "failed to get list of not finished orders")
	}

	return orders, nil
}
func (o *OrderStore) UpdateOrders(ctx context.Context, orders []domain.Order) error {
	if len(orders) == 0 {
		return nil
	}

	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	stmt, err := tx.Prepare(`
		update orders
		set status=$1, accrual=$2
		where order_number=$1
	`)

	if err != nil {
		return errors.Wrap(err, "failed to prepare the update query")
	}

	defer stmt.Close()

	for _, order := range orders {
		if _, err := stmt.Exec(order.Status, order.Accrual, order.OrderNumber); err != nil {
			if err := tx.Rollback(); err != nil {
				return errors.Wrap(err, "unable to rollback")
			}
			return errors.Wrapf(err, "failed to exec query with order %v", order)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit")
	}

	return nil
}

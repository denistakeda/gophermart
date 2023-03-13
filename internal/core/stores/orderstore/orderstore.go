package orderstore

import (
	"context"
	"gophermart/internal/core/domain"

	"github.com/jmoiron/sqlx"
)

type OrderStore struct{}

func New(db *sqlx.DB) *OrderStore {
	return &OrderStore{}
}

func (o *OrderStore) GetOrder(ctx context.Context, orderNumber int) (domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OrderStore) AddNewOrder(ctx context.Context, userID int, orderNumber int) error {
	//TODO implement me
	panic("implement me")
}

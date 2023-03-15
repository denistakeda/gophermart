package withdrawstore

import (
	"context"
	"gophermart/internal/core/domain"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type WithdrawStore struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *WithdrawStore {
	return &WithdrawStore{db: db}
}

func (w *WithdrawStore) GetAllWithdrawals(ctx context.Context, userID int) ([]domain.Withdrawn, error) {
	var withdrawals []domain.Withdrawn
	if err := w.db.SelectContext(ctx, &withdrawals, `
		select * from withdrawals
		where user_id=$1
		order by processed_at
	`, userID); err != nil {
		return withdrawals, errors.Wrapf(err, "failed to get list of withdrawals for user with id %d", userID)
	}

	return withdrawals, nil
}

func (w *WithdrawStore) AddNewWithdrawn(ctx context.Context, orderNumber string, sum float64, userID int) error {
	if _, err := w.db.ExecContext(ctx, `
		insert into withdrawals("order", sum, user_id, processed_at)
		values ($1, $2, $3, $4)
	`, orderNumber, sum, userID, time.Now()); err != nil {
		return errors.Wrapf(err, "failed to insert withdrawn for order '%s' into a database", orderNumber)
	}

	return nil
}

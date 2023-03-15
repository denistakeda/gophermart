package withdrawstore

import (
	"context"
	"gophermart/internal/core/domain"

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
	`, userID); err != nil {
		return withdrawals, errors.Wrapf(err, "failed to get list of withdrawals for user with id %d", userID)
	}

	return withdrawals, nil
}

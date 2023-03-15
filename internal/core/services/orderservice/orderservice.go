package orderservice

import (
	"context"
	"gophermart/internal/core/apperrors"
	"gophermart/internal/core/domain"
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type OrderService struct {
	logger         zerolog.Logger
	orderStore     ports.OrderStore
	withdrawnStore ports.WithdrawnStore
}

func New(
	logService *logging.LoggerService,
	orderStore ports.OrderStore,
	withdrawnStore ports.WithdrawnStore,
) *OrderService {
	return &OrderService{
		logger:         logService.ComponentLogger("OrderService"),
		orderStore:     orderStore,
		withdrawnStore: withdrawnStore,
	}
}

func (o *OrderService) AddOrder(ctx context.Context, user *domain.User, orderNumber string) error {
	if !luhnValid(orderNumber) {
		return errors.Wrap(apperrors.ErrIncorrectOrderFormat, "incorrect order number")
	}

	order, err := o.orderStore.GetOrder(ctx, orderNumber)
	if err != nil && !errors.Is(err, apperrors.ErrNoSuchOrder) {
		return errors.Wrap(err, "failed to create an order")
	}

	if err == nil {
		if user.ID == order.UserID {
			return errors.Wrap(
				apperrors.ErrOrderWasPostedByThisUser,
				"such order was already posted by this user",
			)
		} else {
			return errors.Wrap(
				apperrors.ErrOrderWasPostedByAnotherUser,
				"such order was already posted by another user",
			)
		}
	}

	return o.orderStore.AddNewOrder(ctx, user.ID, orderNumber)
}

func (o *OrderService) GetAllOrders(ctx context.Context, user *domain.User) ([]domain.Order, error) {
	orders, err := o.orderStore.GetAllOrders(ctx, user.ID)
	if err != nil {
		return []domain.Order{}, errors.Wrapf(err, "failed to get all orders for the user %s", user.Login)
	}

	return orders, nil
}

func (o *OrderService) GetUserBalance(ctx context.Context, user *domain.User) (domain.UserBalance, error) {
	var balance domain.UserBalance

	orders, err := o.GetAllOrders(ctx, user)
	if err != nil {
		return balance, err
	}

	withdrawals, err := o.withdrawnStore.GetAllWithdrawals(ctx, user.ID)
	if err != nil {
		return balance, errors.Wrapf(err, "failed to get all withdrawals for the user %s", user.Login)
	}

	var total float64
	for _, order := range orders {
		if order.Status == domain.OrderStatusProcessed {
			total += order.Accrual
		}
	}

	var spent float64
	for _, withdrawn := range withdrawals {
		spent += withdrawn.Sum
	}

	balance.Current = total - spent
	balance.Withdrawn = spent

	return balance, nil
}

func (o *OrderService) Withdraw(ctx context.Context, orderNumber string, sum float64, user *domain.User) error {
	if !luhnValid(orderNumber) {
		return errors.Wrap(apperrors.ErrIncorrectOrderFormat, "incorrect order number")
	}

	balance, err := o.GetUserBalance(ctx, user)
	if err != nil {
		return errors.Wrapf(err, "failed to get balance for user %s", user.Login)
	}

	if balance.Current < sum {
		return errors.Wrapf(apperrors.ErrNotEnoughMoney, "user %s does not have enough money", user.Login)
	}

	if err = o.withdrawnStore.AddNewWithdrawn(ctx, orderNumber, sum, user.ID); err != nil {
		return errors.Wrapf(err, "failed to create a withdrawn for user %s", user.Login)
	}

	return nil
}

package orderservice

import (
	"context"
	"gophermart/internal/core/apperrors"
	"gophermart/internal/core/domain"
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"
	"strconv"

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
	intOrder, err := strconv.Atoi(orderNumber)
	if err != nil {
		return errors.Wrap(apperrors.ErrIncorrectOrderFormat, err.Error())
	}

	if !luhnValid(intOrder) {
		return errors.Wrap(apperrors.ErrIncorrectOrderFormat, "incorrect order number")
	}

	order, err := o.orderStore.GetOrder(ctx, intOrder)
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

	return o.orderStore.AddNewOrder(ctx, user.ID, intOrder)
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

	withdrawns, err := o.withdrawnStore.GetAllWithdrawns(ctx, user.ID)
	if err != nil {
		return balance, errors.Wrapf(err, "failed to get all withdrawns for the user %s", user.Login)
	}

	var total float64
	for _, order := range orders {
		if order.Status == domain.OrderStatusProcessed {
			total += order.Accrual
		}
	}

	var spent float64
	for _, withdrawn := range withdrawns {
		spent += withdrawn.Sum
	}

	balance.Current = total - spent
	balance.Withdrawn = spent

	return balance, nil
}

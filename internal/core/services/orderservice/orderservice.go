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
	logger     zerolog.Logger
	orderStore ports.OrderStore
}

func New(logService *logging.LoggerService, orderStore ports.OrderStore) *OrderService {
	return &OrderService{
		logger:     logService.ComponentLogger("OrderService"),
		orderStore: orderStore,
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

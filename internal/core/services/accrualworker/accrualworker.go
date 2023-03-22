package accrualworker

import (
	"context"
	"gophermart/internal/core/domain"
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"
	"time"

	"github.com/rs/zerolog"
)

var checkInterval = 5 * time.Second

type AccrualWorker struct {
	accrualService ports.AccrualService
	orderStore     ports.OrderStore
	logger         zerolog.Logger
	stopChan       chan struct{}
}

func New(
	accrualService ports.AccrualService,
	orderStore ports.OrderStore,
	logService *logging.LoggerService,
) *AccrualWorker {
	return &AccrualWorker{
		orderStore:     orderStore,
		logger:         logService.ComponentLogger("AccrualWorker"),
		stopChan:       make(chan struct{}),
		accrualService: accrualService,
	}
}

func (a *AccrualWorker) Run() {
	ticker := time.NewTicker(checkInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				a.checkOrders()
			case <-a.stopChan:
				return
			}
		}
	}()
}

func (a *AccrualWorker) Stop() {
	close(a.stopChan)
}

func (a *AccrualWorker) checkOrders() {
	a.logger.Info().Msg("process unfinished orders")
	ctx, cancel := context.WithTimeout(context.Background(), checkInterval)
	defer cancel()

	orders, err := a.orderStore.GetAllNotFinished(ctx)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to get list of orders to check")
		return
	}

	var updatedOrders []domain.Order

	for _, order := range orders {
		resp, err := a.accrualService.CheckAccrual(order.OrderNumber)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to fetch order from accrual system")
			continue
		}

		if resp.Status != string(order.Status) {
			order.Status = domain.OrderStatus(resp.Status)
			order.Accrual = int(resp.Accrual * 100)
			updatedOrders = append(updatedOrders, order)
		}
	}

	if err := a.orderStore.UpdateOrders(ctx, updatedOrders); err != nil {
		a.logger.Error().Err(err).Msg("failed to update orders")
	}
}

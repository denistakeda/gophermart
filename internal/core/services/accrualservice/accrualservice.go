package accrualservice

import (
	"context"
	"encoding/json"
	"fmt"
	"gophermart/internal/core/domain"
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var checkInterval = 5 * time.Second

type AccrualService struct {
	accrualSystemAddress string
	orderStore           ports.OrderStore
	logger               zerolog.Logger
	stopChan             chan struct{}
}

func New(accrualSystemAddress string, orderStore ports.OrderStore, logService *logging.LoggerService) *AccrualService {
	return &AccrualService{
		accrualSystemAddress: accrualSystemAddress,
		orderStore:           orderStore,
		logger:               logService.ComponentLogger("AccrualService"),
		stopChan:             make(chan struct{}),
	}
}

func (a *AccrualService) Run() {
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

func (a *AccrualService) Stop() {
	close(a.stopChan)
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (a *AccrualService) checkOrders() {
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
		resp, err := a.checkAccrual(order)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to fetch order from accrual system")
			continue
		}

		if resp.Status != string(order.Status) {
			order.Status = domain.OrderStatus(resp.Status)
			order.Accrual = resp.Accrual
			updatedOrders = append(updatedOrders, order)
		}
	}

	if err := a.orderStore.UpdateOrders(ctx, updatedOrders); err != nil {
		a.logger.Error().Err(err).Msg("failed to update orders")
	}
}

func (a *AccrualService) checkAccrual(order domain.Order) (AccrualResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", a.accrualSystemAddress, order.OrderNumber))
	if err != nil {
		return AccrualResponse{}, errors.Wrapf(err, "failed to get order '%s' from accrual system", order.OrderNumber)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.logger.Error().Err(err).Msg("failed to close body")
		}
	}()

	var accrualResponse AccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&accrualResponse); err != nil {
		return AccrualResponse{}, errors.Wrapf(err, "failed to decode body for order '%s'", order.OrderNumber)
	}

	return accrualResponse, nil

}

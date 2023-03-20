package accrualservice

import (
	"encoding/json"
	"fmt"
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type AccrualService struct {
	accrualSystemAddress string
	logger               zerolog.Logger
}

func New(accrualSystemAddress string, logService *logging.LoggerService) *AccrualService {
	return &AccrualService{
		accrualSystemAddress: accrualSystemAddress,
		logger:               logService.ComponentLogger("AccrualService"),
	}
}

func (a *AccrualService) CheckAccrual(orderNumber string) (ports.AccrualResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/orders/%s", a.accrualSystemAddress, orderNumber))
	if err != nil {
		return ports.AccrualResponse{}, errors.Wrapf(err, "failed to get order '%s' from accrual system", orderNumber)
	}

	if resp.StatusCode != http.StatusOK {
		return ports.AccrualResponse{}, errors.Errorf("accrual service wrong status: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	var accrualResponse ports.AccrualResponse
	if err := json.NewDecoder(resp.Body).Decode(&accrualResponse); err != nil {
		return ports.AccrualResponse{}, errors.Wrapf(err, "failed to decode body for order '%s'", orderNumber)
	}

	return accrualResponse, nil

}

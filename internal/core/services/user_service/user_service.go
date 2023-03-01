package user_service

import (
	"gophermart/internal/core/services/logging"

	"github.com/rs/zerolog"
)

type UserService struct {
	logger zerolog.Logger
}

func New(logService *logging.LoggerService) *UserService {
	return &UserService{
		logger: logService.ComponentLogger("UserService"),
	}
}

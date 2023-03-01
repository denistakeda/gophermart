package user_api

import (
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserAPI struct {
	logger      zerolog.Logger
	userService ports.UserService
}

func New(logService *logging.LoggerService, userService ports.UserService) *UserAPI {
	return &UserAPI{
		logger:      logService.ComponentLogger("UserAPI"),
		userService: userService,
	}
}

func (api *UserAPI) Register(engine *gin.Engine) {
	engine.Group("/user")
}

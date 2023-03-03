package user_api

import (
	"fmt"
	"gophermart/internal/core/app_errors"
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
	userGroup := engine.Group("/api/user")

	userGroup.POST("/register", api.registerUser)
}

type registerUserBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (api *UserAPI) registerUser(c *gin.Context) {
	var body registerUserBody
	if err := c.ShouldBindJSON(&body); err != nil {
		api.reportError(c, err, http.StatusBadRequest, "invalid body")
		return
	}

	token, err := api.userService.RegisterUser(c, body.Login, body.Password)
	if errors.Is(err, app_errors.ErrLoginIsBusy) {
		api.reportError(c, err, http.StatusConflict, "login is busy")
		return
	} else if err != nil {
		api.reportError(c, err, http.StatusInternalServerError, "server error")
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (api *UserAPI) reportError(c *gin.Context, err error, status int, msg string) {
	api.logger.Error().Err(err).Msg(msg)
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}

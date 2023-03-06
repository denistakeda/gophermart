package userapi

import (
	"fmt"
	"gophermart/internal/core/apperrors"
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var (
	AuthorizationHeaderName = "Authorization"
	UserKey                 = "user"
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

	userGroup.POST("/register", api.registerUserHandler)
	userGroup.POST("/login", api.loginUserHandler)
}

type userBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (api *UserAPI) registerUserHandler(c *gin.Context) {
	var body userBody
	if err := c.ShouldBindJSON(&body); err != nil {
		api.reportError(c, err, http.StatusBadRequest, "invalid body")
		return
	}

	token, err := api.userService.RegisterUser(c, body.Login, body.Password)
	if errors.Is(err, apperrors.ErrLoginIsBusy) {
		api.reportError(c, err, http.StatusConflict, "login is busy")
		return
	} else if err != nil {
		api.reportError(c, err, http.StatusInternalServerError, "server error")
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (api *UserAPI) loginUserHandler(c *gin.Context) {
	var body userBody
	if err := c.ShouldBindJSON(&body); err != nil {
		api.reportError(c, err, http.StatusBadRequest, "invalid body")
		return
	}

	token, err := api.userService.LoginUser(c, body.Login, body.Password)
	if errors.Is(err, apperrors.ErrLoginOrPasswordIncorrect) {
		api.reportError(c, err, http.StatusUnauthorized, "login or password incorrect")
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

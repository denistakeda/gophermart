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
	logger       zerolog.Logger
	userService  ports.UserService
	orderService ports.OrderService
}

func New(
	logService *logging.LoggerService,
	userService ports.UserService,
	orderService ports.OrderService,
) *UserAPI {
	return &UserAPI{
		logger:       logService.ComponentLogger("UserAPI"),
		userService:  userService,
		orderService: orderService,
	}
}

func (api *UserAPI) Register(engine *gin.Engine) {
	userGroup := engine.Group("/api/user")

	userGroup.POST("/register", api.registerUserHandler)
	userGroup.POST("/login", api.loginUserHandler)

	userGroup.POST("/orders", api.AuthMiddleware, api.registerOrderHandler)
	userGroup.GET("/orders", api.AuthMiddleware, api.getOrdersHandler)

	balanceGroup := userGroup.Group("/balance")
	balanceGroup.GET("/", api.AuthMiddleware, api.balanceHandler)
	balanceGroup.POST("/withdraw", api.AuthMiddleware, api.withdrawHandler)

	userGroup.GET("/withdrawals", api.AuthMiddleware, api.withdrawalsHandler)
}

type userBody struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (api *UserAPI) registerUserHandler(c *gin.Context) {
	var body userBody
	if err := c.ShouldBindJSON(&body); err != nil {
		reportError(c, "invalid body", http.StatusBadRequest)
		return
	}

	token, err := api.userService.RegisterUser(c, body.Login, body.Password)
	if errors.Is(err, apperrors.ErrLoginIsBusy) {
		reportError(c, "login is busy", http.StatusConflict)
		return
	} else if err != nil {
		reportError(c, "server error", http.StatusInternalServerError)
		api.logger.Error().Err(err).Msg("failed to register user")
		return
	}

	c.Header(AuthorizationHeaderName, fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (api *UserAPI) loginUserHandler(c *gin.Context) {
	var body userBody
	if err := c.ShouldBindJSON(&body); err != nil {
		reportError(c, "invalid body", http.StatusBadRequest)
		return
	}

	token, err := api.userService.LoginUser(c, body.Login, body.Password)
	if errors.Is(err, apperrors.ErrLoginOrPasswordIncorrect) {
		reportError(c, "login or password incorrect", http.StatusUnauthorized)
		return
	} else if err != nil {
		reportError(c, "server error", http.StatusInternalServerError)
		api.logger.Error().Err(err).Msg("failed to login user")
		return
	}

	c.Header(AuthorizationHeaderName, fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (api *UserAPI) registerOrderHandler(c *gin.Context) {
	user := api.GetUser(c)

	body, err := c.GetRawData()
	if err != nil {
		reportError(c, "order should be a number", http.StatusBadRequest)
		return
	}

	err = api.orderService.AddOrder(c, &user, string(body))
	switch {
	case errors.Is(err, apperrors.ErrOrderWasPostedByThisUser):
		c.JSON(http.StatusOK, gin.H{"success": "order with such number was already posted by this user"})
	case errors.Is(err, apperrors.ErrOrderWasPostedByAnotherUser):
		reportError(c, "order with such number was already posted by another user", http.StatusConflict)
	case errors.Is(err, apperrors.ErrIncorrectOrderFormat):
		reportError(c, "incorrect order format", http.StatusUnprocessableEntity)
	case err != nil:
		reportError(c, "internal server error", http.StatusInternalServerError)
		api.logger.Error().Err(err).Msg("failed to register order")
	default:
		c.JSON(http.StatusAccepted, gin.H{"success": "order was successfully created"})
	}

}

func (api *UserAPI) getOrdersHandler(c *gin.Context) {
	user := api.GetUser(c)

	orders, err := api.orderService.GetAllOrders(c, &user)
	if err != nil {
		reportError(c, "failed to fetch list of orders for a user", http.StatusInternalServerError)
		api.logger.Error().Err(err).Msg("failed to fetch orders")
		return
	}

	if len(orders) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (api *UserAPI) balanceHandler(c *gin.Context) {
	user := api.GetUser(c)

	balance, err := api.orderService.GetUserBalance(c, &user)
	if err != nil {
		reportError(c, "failed to get user balance", http.StatusInternalServerError)
		api.logger.Error().Err(err).Msg("failed to fetch user balance")
		return
	}

	c.JSON(http.StatusOK, balance)
}

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (api *UserAPI) withdrawHandler(c *gin.Context) {
	user := api.GetUser(c)

	var request withdrawRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		reportError(c, "wrong withdraw", http.StatusBadRequest)
		return
	}

	err := api.orderService.Withdraw(c, request.Order, request.Sum, &user)
	switch {
	case errors.Is(err, apperrors.ErrNotEnoughMoney):
		reportError(c, "not enough money", http.StatusPaymentRequired)
	case errors.Is(err, apperrors.ErrIncorrectOrderFormat):
		reportError(c, "incorrect order format", http.StatusUnprocessableEntity)
	case err != nil:
		reportError(c, "internal server error", http.StatusInternalServerError)
		api.logger.Error().Err(err).Msg("failed to fetch withdraw points")
	default:
		c.JSON(http.StatusOK, gin.H{"success": "OK"})
	}
}

func (api *UserAPI) withdrawalsHandler(c *gin.Context) {
	user := api.GetUser(c)

	withdrawals, err := api.orderService.GetAllWithdrawals(c, &user)
	if err != nil {
		reportError(c, "internal server error", http.StatusInternalServerError)
		api.logger.Error().Err(err).Msg("failed to fetch a list of withdrawals")
		return
	}

	if len(withdrawals) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}

func reportError(c *gin.Context, msg string, status int) {
	//api.logger.Error().Err(err).Msg(msg)
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}

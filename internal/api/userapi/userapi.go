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

	userGroup.GET("/balance", api.AuthMiddleware, api.balanceHandler)

	userGroup.POST("/withdraw", api.AuthMiddleware, api.withdrawHandler)
}

type userBody struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	c.Header(AuthorizationHeaderName, fmt.Sprintf("Bearer %s", token))
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

	c.Header(AuthorizationHeaderName, fmt.Sprintf("Bearer %s", token))
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (api *UserAPI) registerOrderHandler(c *gin.Context) {
	user := api.GetUser(c)

	body, err := c.GetRawData()
	if err != nil {
		api.reportError(c, err, http.StatusBadRequest, "order should be a number")
		return
	}

	err = api.orderService.AddOrder(c, &user, string(body))
	switch {
	case errors.Is(err, apperrors.ErrOrderWasPostedByThisUser):
		c.JSON(http.StatusOK, gin.H{"success": "order with such number was already posted by this user"})
	case errors.Is(err, apperrors.ErrOrderWasPostedByAnotherUser):
		api.reportError(c, err, http.StatusConflict, "order with such number was already posted by another user")
	case errors.Is(err, apperrors.ErrIncorrectOrderFormat):
		api.reportError(c, err, http.StatusUnprocessableEntity, "incorrect order format")
	case err != nil:
		api.reportError(c, err, http.StatusInternalServerError, "internal server error")
	default:
		c.JSON(http.StatusAccepted, gin.H{"success": "order was successfully created"})
	}

}

func (api *UserAPI) getOrdersHandler(c *gin.Context) {
	user := api.GetUser(c)

	orders, err := api.orderService.GetAllOrders(c, &user)
	if err != nil {
		api.reportError(c, err, http.StatusInternalServerError, "failed to fetch list of orders for a user")
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
		api.reportError(c, err, http.StatusInternalServerError, "failed to get user balance")
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
		api.reportError(c, err, http.StatusBadRequest, "wrong withdraw")
		return
	}

	err := api.orderService.Withdraw(c, request.Order, request.Sum, user.ID)
	switch {
	case errors.Is(err, apperrors.ErrNotEnoughMoney):
		api.reportError(c, err, http.StatusPaymentRequired, "not enough money")
	case errors.Is(err, apperrors.ErrIncorrectOrderFormat):
		api.reportError(c, err, http.StatusUnprocessableEntity, "incorrect order format")
	case err != nil:
		api.reportError(c, err, http.StatusInternalServerError, "internal server error")
	default:
		c.JSON(http.StatusOK, gin.H{"success": "OK"})
	}
}

func (api *UserAPI) reportError(c *gin.Context, err error, status int, msg string) {
	api.logger.Error().Err(err).Msg(msg)
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}

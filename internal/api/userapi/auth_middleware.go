package userapi

import (
	"gophermart/internal/core/apperrors"
	"gophermart/internal/core/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (api *UserAPI) AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader(AuthorizationHeaderName)
	if authHeader == "" {
		msg := "missing auth header"
		api.reportError(c, errors.New(msg), http.StatusUnauthorized, msg)
		return
	}

	split := strings.Split(authHeader, " ")
	if len(split) < 2 || split[0] != "Bearer" {
		msg := "auth header incorrect"
		api.reportError(c, errors.New(msg), http.StatusUnauthorized, msg)
		return
	}

	user, err := api.userService.AuthenticateUser(c, split[1])
	if errors.Is(err, apperrors.ErrAuthFailed) {
		msg := "auth incorrect"
		api.reportError(c, errors.New(msg), http.StatusUnauthorized, msg)
		return
	} else if err != nil {
		msg := "failed to authenticate user"
		api.reportError(c, errors.New(msg), http.StatusInternalServerError, msg)
		return
	}

	c.Set(UserKey, user)

	c.Next()
}

func (api *UserAPI) GetUser(c *gin.Context) domain.User {
	user := c.MustGet(UserKey)
	return user.(domain.User)
}

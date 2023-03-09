package userapi

import (
	"gophermart/internal/core/apperrors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (api *UserAPI) authMiddleware(c *gin.Context) {
	authHeader := c.GetHeader(AuthorizationHeaderName)
	if authHeader == "" {
		msg := "missing auth header"
		api.reportError(c, errors.New(msg), http.StatusUnauthorized, msg)
		return
	}

	splited := strings.Split(authHeader, " ")
	if len(splited) < 2 || splited[0] != "Bearer" {
		msg := "auth header incorrect"
		api.reportError(c, errors.New(msg), http.StatusUnauthorized, msg)
		return
	}

	user, err := api.userService.AuthenticateUser(c, splited[1])
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

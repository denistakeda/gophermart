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
		reportError(c, "missing auth header", http.StatusUnauthorized)
		return
	}

	split := strings.Split(authHeader, " ")
	if len(split) < 2 || split[0] != "Bearer" {
		reportError(c, "auth header incorrect", http.StatusUnauthorized)
		return
	}

	user, err := api.userService.AuthenticateUser(c, split[1])
	if errors.Is(err, apperrors.ErrAuthFailed) {
		reportError(c, "auth incorrect", http.StatusUnauthorized)
		return
	} else if err != nil {
		msg := "failed to authenticate user"
		reportError(c, msg, http.StatusInternalServerError)
		api.logger.Error().Err(err).Msg(msg)
		return
	}

	c.Set(UserKey, user)

	c.Next()
}

func (api *UserAPI) GetUser(c *gin.Context) domain.User {
	user := c.MustGet(UserKey)
	return user.(domain.User)
}

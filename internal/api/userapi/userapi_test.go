package userapi

import (
	"bytes"
	"encoding/json"
	"gophermart/internal/core/apperrors"
	"gophermart/internal/core/services/logging"
	mocks "gophermart/mocks/core/ports"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterNewUser(t *testing.T) {
	type want struct {
		status               int
		shouldHaveAuthHeader bool
	}
	type serviceCall struct {
		args    []any
		returns []any
		times   int
	}
	tests := []struct {
		name        string
		requestBody []byte
		serviceCall *serviceCall

		want want
	}{
		{
			name:        "empty body",
			requestBody: []byte("hello there"),
			want: want{
				status:               http.StatusBadRequest,
				shouldHaveAuthHeader: false,
			},
		},
		{
			name:        "user already exists",
			requestBody: makeRegisterUserBody(t, "user", "password"),
			serviceCall: &serviceCall{
				args:    []interface{}{gomock.Any(), "user", "password"},
				returns: []interface{}{"", errors.Wrap(apperrors.ErrLoginIsBusy, "test error")},
				times:   1,
			},
			want: want{
				status:               http.StatusConflict,
				shouldHaveAuthHeader: false,
			},
		},
		{
			name:        "internal error",
			requestBody: makeRegisterUserBody(t, "user", "password"),
			serviceCall: &serviceCall{
				args:    []interface{}{gomock.Any(), "user", "password"},
				returns: []interface{}{"", errors.New("test error")},
				times:   1,
			},
			want: want{
				status:               http.StatusInternalServerError,
				shouldHaveAuthHeader: false,
			},
		},
		{
			name:        "success case",
			requestBody: makeRegisterUserBody(t, "user", "password"),
			serviceCall: &serviceCall{
				args:    []interface{}{gomock.Any(), "user", "password"},
				returns: []interface{}{"token", nil},
				times:   1,
			},
			want: want{
				status:               http.StatusOK,
				shouldHaveAuthHeader: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			userService := mocks.NewMockUserService(ctrl)
			router := gin.New()
			logService := logging.New()
			userAPI := New(logService, userService)
			userAPI.Register(router)

			if tt.serviceCall != nil {
				userService.EXPECT().
					RegisterUser(tt.serviceCall.args[0], tt.serviceCall.args[1], tt.serviceCall.args[2]).
					Return(tt.serviceCall.returns[0], tt.serviceCall.returns[1]).
					Times(tt.serviceCall.times)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.status, w.Code)

			header := w.Header().Get("Authorization")
			if tt.want.shouldHaveAuthHeader {
				assert.NotEmptyf(t, header, "authorization header should be presented")
				assert.True(
					t,
					strings.Contains(header, "Bearer"),
					"authorization header should start with 'Bearer'",
				)
			} else {
				assert.Empty(t, header)
			}
		})
	}
}

func makeRegisterUserBody(t *testing.T, login, password string) []byte {
	body := registerUserBody{
		Login:    login,
		Password: password,
	}
	res, err := json.Marshal(body)
	require.NoError(t, err)
	return res
}

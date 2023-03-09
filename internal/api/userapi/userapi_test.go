package userapi

import (
	"bytes"
	"encoding/json"
	"gophermart/internal/core/apperrors"
	"gophermart/internal/core/domain"
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
			name:        "incorrect body",
			requestBody: []byte("hello there"),
			want: want{
				status:               http.StatusBadRequest,
				shouldHaveAuthHeader: false,
			},
		},
		{
			name:        "user already exists",
			requestBody: makeUserBody(t, "user", "password"),
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
			requestBody: makeUserBody(t, "user", "password"),
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
			requestBody: makeUserBody(t, "user", "password"),
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

func TestLoginUser(t *testing.T) {
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
			name:        "incorrect body",
			requestBody: []byte("hello there"),
			want: want{
				status:               http.StatusBadRequest,
				shouldHaveAuthHeader: false,
			},
		},
		{
			name:        "incorrect login or password",
			requestBody: makeUserBody(t, "user", "password"),
			serviceCall: &serviceCall{
				args:    []interface{}{gomock.Any(), "user", "password"},
				returns: []interface{}{"", errors.Wrap(apperrors.ErrLoginOrPasswordIncorrect, "test error")},
				times:   1,
			},
			want: want{
				status:               http.StatusUnauthorized,
				shouldHaveAuthHeader: false,
			},
		},
		{
			name:        "internal error",
			requestBody: makeUserBody(t, "user", "password"),
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
			requestBody: makeUserBody(t, "user", "password"),
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
					LoginUser(tt.serviceCall.args[0], tt.serviceCall.args[1], tt.serviceCall.args[2]).
					Return(tt.serviceCall.returns[0], tt.serviceCall.returns[1]).
					Times(tt.serviceCall.times)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewReader(tt.requestBody))
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

func TestRegisterOrder(t *testing.T) {
	type want struct {
		status int
	}
	type serviceCall struct {
		args    []any
		returns []any
		times   int
	}
	tests := []struct {
		name                 string
		requestBody          []byte
		authHeader           string
		authenticateUserCall *serviceCall

		want want
	}{
		{
			name:        "missing auth header",
			requestBody: []byte("test body"),
			authHeader:  "",
			// authenticateUserCall: &serviceCall{
			// 	args: []any{gomock.Any(), }
			// },
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name:        "incorrect auth header",
			requestBody: []byte("test body"),
			authHeader:  "Onetwothree",
			// authenticateUserCall: &serviceCall{
			// 	args: []any{gomock.Any(), }
			// },
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name:        "missing auth token",
			requestBody: []byte("test body"),
			authHeader:  "Bearer",
			// authenticateUserCall: &serviceCall{
			// 	args: []any{gomock.Any(), }
			// },
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name:        "authentication failed",
			requestBody: []byte("test body"),
			authHeader:  "Bearer authtoken",
			authenticateUserCall: &serviceCall{
				args:    []any{gomock.Any(), "authtoken"},
				returns: []any{domain.User{}, errors.Wrap(apperrors.ErrAuthFailed, "test error")},
				times:   1,
			},
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name:        "unknown service error",
			requestBody: []byte("test body"),
			authHeader:  "Bearer authtoken",
			authenticateUserCall: &serviceCall{
				args:    []any{gomock.Any(), "authtoken"},
				returns: []any{domain.User{}, errors.New("test error")},
				times:   1,
			},
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		// TODO: remove this test
		{
			name:        "successful authentication",
			requestBody: []byte("test body"),
			authHeader:  "Bearer authtoken",
			authenticateUserCall: &serviceCall{
				args:    []any{gomock.Any(), "authtoken"},
				returns: []any{domain.User{Login: "user", Password: "password"}, nil},
				times:   1,
			},
			want: want{
				status: http.StatusOK,
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

			if tt.authenticateUserCall != nil {
				userService.EXPECT().
					AuthenticateUser(tt.authenticateUserCall.args[0], tt.authenticateUserCall.args[1]).
					Return(tt.authenticateUserCall.returns[0], tt.authenticateUserCall.returns[1]).
					Times(tt.authenticateUserCall.times)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/orders", bytes.NewReader(tt.requestBody))

			req.Header.Set("Content-Type", "text/plain")

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.status, w.Code)
		})
	}
}

func makeUserBody(t *testing.T, login, password string) []byte {
	body := userBody{
		Login:    login,
		Password: password,
	}
	res, err := json.Marshal(body)
	require.NoError(t, err)
	return res
}

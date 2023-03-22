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
			apiTest := NewAPITest(t)

			if tt.serviceCall != nil {
				apiTest.UserService.EXPECT().
					RegisterUser(tt.serviceCall.args[0], tt.serviceCall.args[1], tt.serviceCall.args[2]).
					Return(tt.serviceCall.returns[0], tt.serviceCall.returns[1]).
					Times(tt.serviceCall.times)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/register", bytes.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			apiTest.Router.ServeHTTP(w, req)

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
			apiTest := NewAPITest(t)

			if tt.serviceCall != nil {
				apiTest.UserService.EXPECT().
					LoginUser(tt.serviceCall.args[0], tt.serviceCall.args[1], tt.serviceCall.args[2]).
					Return(tt.serviceCall.returns[0], tt.serviceCall.returns[1]).
					Times(tt.serviceCall.times)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			apiTest.Router.ServeHTTP(w, req)

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

func TestAuthMiddleware(t *testing.T) {
	type want struct {
		status int
	}
	type authenticateUserCall struct {
		args    []any
		returns []any
		times   int
	}
	tests := []struct {
		name                 string
		requestBody          []byte
		authHeader           string
		authenticateUserCall *authenticateUserCall

		want want
	}{
		{
			name:        "missing auth header",
			requestBody: []byte("test body"),
			authHeader:  "",
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name:        "incorrect auth header",
			requestBody: []byte("test body"),
			authHeader:  "Onetwothree",
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name:        "missing auth token",
			requestBody: []byte("test body"),
			authHeader:  "Bearer",
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name:        "authentication failed",
			requestBody: []byte("test body"),
			authHeader:  "Bearer authtoken",
			authenticateUserCall: &authenticateUserCall{
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
			authenticateUserCall: &authenticateUserCall{
				args:    []any{gomock.Any(), "authtoken"},
				returns: []any{domain.User{}, errors.New("test error")},
				times:   1,
			},
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		{
			name:        "successfull authorization",
			requestBody: []byte("text-order-number"),
			authHeader:  "Bearer authtoken",
			authenticateUserCall: &authenticateUserCall{
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
			apiTest := NewAPITest(t)
			apiTest.Router.GET("/authmiddleware", apiTest.UserAPI.AuthMiddleware, func(c *gin.Context) {
				val, exists := c.Get(UserKey)
				assert.True(t, exists, "middleware should assign a user")
				_, ok := val.(domain.User)
				assert.True(t, ok, "middleware should assign a value of type user")
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			if tt.authenticateUserCall != nil {
				apiTest.UserService.EXPECT().
					AuthenticateUser(tt.authenticateUserCall.args[0], tt.authenticateUserCall.args[1]).
					Return(tt.authenticateUserCall.returns[0], tt.authenticateUserCall.returns[1]).
					Times(tt.authenticateUserCall.times)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/authmiddleware", bytes.NewReader(tt.requestBody))

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			apiTest.Router.ServeHTTP(w, req)

			assert.Equal(t, tt.want.status, w.Code)
		})
	}
}

func TestRegisterOrderHandler(t *testing.T) {
	orderNum := "1234 5678 1234 5678"

	tests := []struct {
		name                string
		user                domain.User
		requestBody         string
		addOrderErrorReturn error

		wantStatus int
	}{
		{
			name:                "already aploaded order",
			user:                domain.User{},
			requestBody:         orderNum,
			addOrderErrorReturn: apperrors.ErrOrderWasPostedByThisUser,
			wantStatus:          http.StatusOK,
		},
		{
			name:                "order was posted by another user",
			user:                domain.User{},
			requestBody:         orderNum,
			addOrderErrorReturn: apperrors.ErrOrderWasPostedByAnotherUser,
			wantStatus:          http.StatusConflict,
		},
		{
			name:                "incorrect order format",
			user:                domain.User{},
			requestBody:         orderNum,
			addOrderErrorReturn: apperrors.ErrIncorrectOrderFormat,
			wantStatus:          http.StatusUnprocessableEntity,
		},
		{
			name:                "internal server error",
			user:                domain.User{},
			requestBody:         orderNum,
			addOrderErrorReturn: errors.New("unknown error"),
			wantStatus:          http.StatusInternalServerError,
		},
		{
			name:                "success case",
			user:                domain.User{},
			requestBody:         orderNum,
			addOrderErrorReturn: nil,
			wantStatus:          http.StatusAccepted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiTest := NewAPITest(t).AuthenticateWithUser(tt.user)

			apiTest.OrderService.EXPECT().
				AddOrder(gomock.Any(), &tt.user, tt.requestBody).
				Return(tt.addOrderErrorReturn).
				Times(1)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/user/orders", bytes.NewReader([]byte(tt.requestBody)))

			req.Header.Set("Authorization", "Bearer authtoken")

			apiTest.Router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// -- Test helpers --

type APITest struct {
	Router       *gin.Engine
	UserService  *mocks.MockUserService
	OrderService *mocks.MockOrderService
	UserAPI      *UserAPI
	LogService   *logging.LoggerService
}

func NewAPITest(t *testing.T) *APITest {
	ctrl := gomock.NewController(t)
	userService := mocks.NewMockUserService(ctrl)
	router := gin.New()
	logService := logging.New()
	orderService := mocks.NewMockOrderService(ctrl)
	userAPI := New(logService, userService, orderService)

	userAPI.Register(router)

	return &APITest{
		Router:       router,
		UserService:  userService,
		OrderService: orderService,
		UserAPI:      userAPI,
		LogService:   logService,
	}
}

func (apiTest *APITest) AuthenticateWithUser(user domain.User) *APITest {
	apiTest.UserService.EXPECT().
		AuthenticateUser(gomock.Any(), gomock.Any()).
		Return(user, nil).
		Times(1)

	return apiTest
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

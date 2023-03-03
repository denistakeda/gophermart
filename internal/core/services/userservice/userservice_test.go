package userservice

import (
	"context"
	"gophermart/internal/core/apperrors"
	"gophermart/internal/core/services/logging"
	mock "gophermart/mocks/core/ports"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestUserService_RegisterUser(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	type want struct {
		errorIs error
		token   bool
	}
	type storeCall struct {
		args    []interface{}
		returns error
		times   int
	}
	tests := []struct {
		name      string
		args      args
		storeCall *storeCall
		want      want
	}{
		{
			name: "empty login",
			args: args{
				login:    "",
				password: "password",
			},
			want: want{
				errorIs: apperrors.ErrLoginIsEmpty,
				token:   false,
			},
		},
		{
			name: "empty password",
			args: args{
				login:    "login",
				password: "",
			},
			want: want{
				errorIs: apperrors.ErrPasswordIsEmpty,
				token:   false,
			},
		},
		{
			name: "busy login",
			args: args{
				login:    "login",
				password: "password",
			},
			storeCall: &storeCall{
				args:    []interface{}{gomock.Any(), "login", gomock.Any()},
				returns: errors.New("test error"),
				times:   1,
			},
			want: want{
				errorIs: apperrors.ErrLoginIsBusy,
				token:   false,
			},
		},
		{
			name: "success case",
			args: args{
				login:    "login",
				password: "password",
			},
			storeCall: &storeCall{
				args:    []interface{}{gomock.Any(), "login", gomock.Any()},
				returns: nil,
				times:   1,
			},
			want: want{
				errorIs: nil,
				token:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logService := logging.New()
			ctrl := gomock.NewController(t)
			userStore := mock.NewMockUserStore(ctrl)
			userService := New("secret", logService, userStore)

			if tt.storeCall != nil {
				userStore.
					EXPECT().
					AddNewUser(tt.storeCall.args[0], tt.storeCall.args[1], tt.storeCall.args[2]).
					Return(tt.storeCall.returns).
					Times(tt.storeCall.times)
			}

			token, err := userService.RegisterUser(context.Background(), tt.args.login, tt.args.password)
			if tt.want.errorIs != nil {
				assert.ErrorIs(t, err, tt.want.errorIs)
			} else {
				assert.NoError(t, err)
			}
			if tt.want.token {
				assert.NotEmptyf(t, token, "token shouldn't be empty")
			}
		})
	}
}

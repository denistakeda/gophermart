package apperrors

import "github.com/pkg/errors"

var (
	ErrLoginIsBusy              = errors.New("login is busy")
	ErrLoginIsEmpty             = errors.New("login is empty")
	ErrPasswordIsEmpty          = errors.New("password is empty")
	ErrLoginOrPasswordIncorrect = errors.New("login or password incorrect")
	ErrAuthFailed               = errors.New("authentication has failed")
)

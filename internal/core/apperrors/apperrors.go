package apperrors

import "github.com/pkg/errors"

var (
	ErrLoginIsBusy     = errors.New("login is busy")
	ErrLoginIsEmpty    = errors.New("login is empty")
	ErrPasswordIsEmpty = errors.New("password is empty")
)

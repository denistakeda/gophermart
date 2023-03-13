package apperrors

import "github.com/pkg/errors"

var (
	ErrLoginIsBusy              = errors.New("login is busy")
	ErrLoginIsEmpty             = errors.New("login is empty")
	ErrPasswordIsEmpty          = errors.New("password is empty")
	ErrLoginOrPasswordIncorrect = errors.New("login or password incorrect")
	ErrAuthFailed               = errors.New("authentication has failed")

	ErrOrderWasPostedByThisUser    = errors.New("the order with such number was already posted by this user")
	ErrOrderWasPostedByAnotherUser = errors.New("the order with such number was already posted by another user")
	ErrIncorrectOrderFormat        = errors.New("incorrect order format")
	ErrNoSuchOrder                 = errors.New("no such order in the database")
)

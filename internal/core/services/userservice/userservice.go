package userservice

import (
	"context"
	"gophermart/internal/core/apperrors"
	"gophermart/internal/core/ports"
	"gophermart/internal/core/services/logging"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	logger    zerolog.Logger
	userStore ports.UserStore
	secret    string
}

func New(secret string, logService *logging.LoggerService, userStore ports.UserStore) *UserService {
	return &UserService{
		logger:    logService.ComponentLogger("UserService"),
		userStore: userStore,
		secret:    secret,
	}
}

func (u *UserService) RegisterUser(ctx context.Context, login, password string) (string, error) {
	if login == "" {
		return "", apperrors.ErrLoginIsEmpty
	}

	if password == "" {
		return "", apperrors.ErrPasswordIsEmpty
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate password hash")
	}

	if err := u.userStore.AddNewUser(ctx, login, string(passwordHash)); err != nil {
		u.logger.Error().Err(err).Msg("failed to add a new user")
		return "", errors.Wrapf(apperrors.ErrLoginIsBusy, "login '%s' is busy", login)
	}

	return u.generateJWT(login)
}

func (u *UserService) LoginUser(ctx context.Context, login, password string) (string, error) {
	if login == "" {
		return "", apperrors.ErrLoginIsEmpty
	}

	if password == "" {
		return "", apperrors.ErrPasswordIsEmpty
	}

	user, err := u.userStore.GetUser(ctx, login)
	if err != nil {
		return "", errors.Wrap(
			apperrors.ErrLoginOrPasswordIncorrect,
			"user with such login/password does not exist",
		)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.Wrap(
			apperrors.ErrLoginOrPasswordIncorrect,
			"user with such login/password does not exist",
		)
	}

	return u.generateJWT(login)
}

func (u *UserService) generateJWT(login string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":      login,
		"created_at": time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(u.secret))
	if err != nil {
		return "", errors.Wrap(err, "failed to create a token")
	}

	return tokenString, nil
}

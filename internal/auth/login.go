package auth

import (
	"errors"
)

var (
	// ErrInvalidPassword is returned when the password is invalid.
	ErrInvalidPassword = errors.New("invalid password")
)

type LoginCommand struct {
	Email    string
	Password string
}

type LoginHandler interface {
	Handle(command LoginCommand) (token Token, err error)
}

type loginHandler struct {
	repo            Repository
	tokenGenerator  TokenHelper[User]
	passwordChecker PasswordChecker
}

func (h *loginHandler) Handle(command LoginCommand) (token Token, err error) {
	user, err := h.repo.FindByEmail(command.Email)
	if err != nil {
		return "", err
	}
	if !h.passwordChecker.Check(command.Password, user.Password) {
		return "", ErrInvalidPassword
	}
	return h.tokenGenerator.Generate(*user), nil
}

func NewLoginHandler(
	repo Repository,
	tokenGenerator TokenHelper[User],
	passChecker PasswordChecker,
) LoginHandler {
	return &loginHandler{repo, tokenGenerator, passChecker}
}

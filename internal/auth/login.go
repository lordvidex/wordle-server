package auth

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/wordle-wf/internal/game"
)

var (
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidPassword is returned when the password is invalid.
	ErrInvalidPassword = errors.New("invalid password")
)

type LoginCommand struct {
	Email    string
	Password string
}

type PlayerWithToken struct {
	Player *game.Player
	Token  Token
}

type LoginHandler interface {
	Handle(command LoginCommand) (token *PlayerWithToken, err error)
}

type loginHandler struct {
	repo            Repository
	tokenGenerator  TokenHelper
	passwordChecker PasswordHelper
}

func (h *loginHandler) Handle(command LoginCommand) (result *PlayerWithToken, err error) {
	user, err := h.repo.FindByEmail(command.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if !h.passwordChecker.Validate(command.Password, user.Password) {
		return nil, ErrInvalidPassword
	}
	token, err := h.tokenGenerator.Generate(user)
	if err != nil {
		return nil, err
	}
	return &PlayerWithToken{user, token}, nil
}

func NewLoginHandler(
	repo Repository,
	tokenGenerator TokenHelper,
	passChecker PasswordHelper,
) LoginHandler {
	return &loginHandler{repo, tokenGenerator, passChecker}
}

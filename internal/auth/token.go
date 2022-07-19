package auth

import (
	"errors"
	"github.com/lordvidex/wordle-wf/internal/game"
)

// Token represents a token used to authenticate a user.
type Token string

var (
	// ErrInvalidToken is returned when the token is invalid.
	ErrInvalidToken = errors.New("invalid token")
)

type GetUserByTokenQueryHandler interface {
	Handle(token Token) (*game.Player, error)
}

type userTokenDecoder struct {
	tokenHelper TokenHelper
}

func (d *userTokenDecoder) Handle(token Token) (*game.Player, error) {
	var payload game.Player
	err := d.tokenHelper.Decode(token, &payload)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return &payload, nil
}

func NewUserTokenDecoder(tokenHelper TokenHelper) GetUserByTokenQueryHandler {
	return &userTokenDecoder{tokenHelper}
}

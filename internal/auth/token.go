package auth

import "errors"

// Token represents a token used to authenticate a user.
type Token string

var (
	// ErrInvalidToken is returned when the token is invalid.
	ErrInvalidToken = errors.New("invalid token")
)

type QueryGetUserByToken interface {
	Handle(token Token) (*User, error)
}

type userTokenDecoder struct {
	tokenHelper TokenHelper
}

func (d *userTokenDecoder) Handle(token Token) (*User, error) {
	user, err := d.tokenHelper.Decode(token)
	if err != nil {
		return nil, err
	}
	switch user.(type) {
	case *User:
		return user.(*User), nil
	case User:
		instance := user.(User)
		return &instance, nil
	default:
		return nil, ErrInvalidToken
	}
}

func NewUserTokenDecoder(tokenHelper TokenHelper) QueryGetUserByToken {
	return &userTokenDecoder{tokenHelper}
}

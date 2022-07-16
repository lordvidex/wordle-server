package auth

// Token represents a token used to authenticate a user.
type Token string

type QueryGetUserByToken interface {
	Handle(token Token) (*User, error)
}

type userTokenDecoder struct {
	tokenHelper TokenHelper[User]
}

func (d *userTokenDecoder) Handle(token Token) (*User, error) {
	user, err := d.tokenHelper.Decode(token)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewUserTokenDecoder(tokenHelper TokenHelper[User]) QueryGetUserByToken {
	return &userTokenDecoder{tokenHelper}
}

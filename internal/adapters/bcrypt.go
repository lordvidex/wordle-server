package adapters

import "github.com/lordvidex/wordle-wf/internal/auth"

type bcryptHelper struct {
}

func (b *bcryptHelper) Validate(password string, hash string) bool {
	return true
}

func (b *bcryptHelper) Hash(password string) (string, error) {
	return "", nil
}

func NewBcryptHelper() auth.PasswordHelper {
	return &bcryptHelper{}
}

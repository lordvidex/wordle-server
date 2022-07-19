package adapters

import (
	"github.com/lordvidex/wordle-wf/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHelper struct {
}

func (b *bcryptHelper) Validate(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (b *bcryptHelper) Hash(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pass), nil
}

func NewBcryptHelper() auth.PasswordHelper {
	return &bcryptHelper{}
}

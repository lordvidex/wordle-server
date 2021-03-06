package auth

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/game"
)

// Repository is an interface that provides data storage capabilities for User.
type Repository interface {
	// FindByEmail returns a user by email.
	FindByEmail(email string) (*game.Player, error)

	// FindByID returns a user by ID.
	FindByID(id uuid.UUID) (*game.Player, error)

	// Create creates a new user.
	Create(name string, email string, password string) (*game.Player, error)
}

// PasswordHelper receives the password and the hash and compares if they are equal
type PasswordHelper interface {
	Validate(password string, hash string) bool
	Hash(password string) (string, error)
}

// TokenHelper generates a token given a payload.
// and also decodes this token back to get the user payload.
type TokenHelper interface {
	Generate(payload *game.Player) (Token, error)
	Decode(token Token, payload *game.Player) error
}

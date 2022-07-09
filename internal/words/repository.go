package words

import "github.com/google/uuid"

type Repository interface {
	Add(word Word) error
	Find(id uuid.UUID) (*Word, error)
	FindAll() ([]Word, error)
}

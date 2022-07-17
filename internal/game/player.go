package game

import (
	"github.com/google/uuid"
)

type Player struct {
	ID        uuid.UUID
	Name      string
	Points    int64
	Email     string
	Password  string `json:"-"`
	IsDeleted bool
}

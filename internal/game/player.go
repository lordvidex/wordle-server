package game

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/auth"
)

type Player struct {
	ID   uuid.UUID
	User *auth.User
	Name string
}

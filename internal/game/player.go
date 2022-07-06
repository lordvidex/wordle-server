package game

import "github.com/google/uuid"

type Player struct {
	ID   uuid.UUID // TODO: change string to UUID
	Name string
}

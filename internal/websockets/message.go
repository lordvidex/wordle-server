package websockets

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/game"
)

// WSGameDto is the game data payload sent to websockets clients
type WSGameDto struct {
	ID       *uuid.UUID      `json:"id"`
	Settings *game.Settings  `json:"settings"`
	Sessions []*game.Session `json:"sessions"`
}

type WSPayload struct {
	Event game.Event  `json:"event"`
	Data  interface{} `json:"data"`
}

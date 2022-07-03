package game

import (
	"github.com/google/uuid"
	"time"
)

type Game struct {
	ID          uuid.UUID
	Word        string
	Home        *Player
	Away        *Player
	HomeSession *Session
	AwaySession *Session
	Settings
}

type Session struct {
	StartTime time.Time
	EndTime   time.Time
	Guesses   []string
}

type Settings struct {
}

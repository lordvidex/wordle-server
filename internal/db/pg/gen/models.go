// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package pg

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Game struct {
	ID     uuid.UUID
	WordID uuid.NullUUID
}

type GameSession struct {
	ID       uuid.UUID
	GameID   uuid.NullUUID
	PlayerID uuid.NullUUID
}

type GameSessionGuess struct {
	ID            uuid.UUID
	GameSessionID uuid.NullUUID
	WordID        uuid.NullUUID
}

type Player struct {
	ID   uuid.UUID
	Name string
}

type Word struct {
	ID         uuid.UUID
	TimePlayed time.Time
	Letters    pgtype.JSON
}

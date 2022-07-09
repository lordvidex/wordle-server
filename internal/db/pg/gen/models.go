// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package pg

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Game struct {
	ID     uuid.UUID
	WordID uuid.NullUUID
}

type GamePlayer struct {
	ID   uuid.UUID
	Name string
}

type GameSession struct {
	ID       uuid.UUID
	GameID   uuid.UUID
	PlayerID uuid.UUID
}

type GameSessionGuess struct {
	ID            uuid.UUID
	GameSessionID uuid.NullUUID
	WordID        uuid.NullUUID
}

type GameSetting struct {
	ID                       uuid.UUID
	GameID                   uuid.NullUUID
	WordLength               sql.NullInt16
	Trials                   sql.NullInt16
	PlayerCount              sql.NullInt16
	HasAnalytics             sql.NullBool
	ShouldRecordTime         sql.NullBool
	CanViewOpponentsSessions sql.NullBool
}

type Word struct {
	ID         uuid.UUID
	TimePlayed time.Time
	Letters    pgtype.JSON
}

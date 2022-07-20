package game

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
)

type NullInt64 struct {
	sql.NullInt64
}

func (nullInt64 *NullInt64) MarshalJSON() ([]byte, error) {
	if !nullInt64.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nullInt64.Int64)
}

func (nullInt64 *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nullInt64.Int64)
	nullInt64.Valid = (err == nil && nullInt64.Int64 > 0)
	return err
}

// Repository is the interface for storing game properties and state
type Repository interface {
	// Create creates a new game and the settings at the same time
	// by default, the game appends to all related tables
	// returns the game instance back newly from the database
	Create(game *Game) (*Game, error)

	// UpdateSettings updates the settings of a game that is yet to be started
	UpdateSettings(settings *Settings, gameID string) error

	// FindByID returns a game by its ID only
	FindByID(gameId uuid.UUID) (*Game, error)

	// UpdateGameResult updates the result of a game for each user in a game
	UpdateGameResult(
		gameID uuid.UUID,
		playerID uuid.UUID,
		playerName string,
		points func(position int) int,
		wordsPlayed []*words.Word,
	) error

	// Delete deletes a game from the database with CASCADE on Settings, Word and Sessions
	Delete(gameId uuid.UUID) error
}

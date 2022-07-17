package game

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
)

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

package game

import (
	"time"

	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
)

const (
	// MaxDuration is the maximum duration a game can last
	MaxDuration = time.Hour
)

type Game struct {
	// ID is the local specific identifier for a game in the local
	ID uuid.UUID

	// InviteID is the id of the lobby used to play this game
	InviteID string

	// Word is the correct word that should be guessed
	Word words.Word

	// Sessions represent each player's game session
	Sessions []*Session

	// Settings represent the rules of the game as set by the room owner
	Settings Settings

	// the number of players that started the game, default is 1 - the creator
	PlayerCount int

	// StartTime is the time the game started,
	StartTime time.Time

	// EndTime is the time the game ended,
	// when the value is nil, this means the game has not ended
	EndTime *time.Time

	CreatorID uuid.UUID
}

func NewGame() *Game {
	return &Game{
		ID:        uuid.New(),
		StartTime: time.Now(),
	}
}

// HasEnded Game if the Game.EndTime is set OR if the game has been active for an hour
// Ended games do not receive rewards after completed Sessions and penalties are applied
// to all sessions immediately after Game has ended.
// or if they have guessed the word correctly
func (g *Game) HasEnded() bool {
	return g.EndTime != nil && g.EndTime.After(g.StartTime.Add(MaxDuration))
}

type Session struct {
	Player  *Player
	Guesses []*words.Word
}

func NewSession(player *Player) *Session {
	return &Session{
		Player: player,
	}
}

type Settings struct {
	WordLength            int  `json:"word_length"`
	Trials                int  `json:"trials"`
	MaxPlayerCount        int  `json:"max_player_count"`
	Analytics             bool `json:"has_analytics"`
	RecordTime            bool `json:"should_record_time"`
	ViewOpponentsSessions bool `json:"can_view_opponents_sessions"`
}

// NewSettings creates a new Settings struct
// with default values
func NewSettings(maxPlayerCount int) Settings {
	return Settings{
		WordLength:            5,
		Trials:                6,
		MaxPlayerCount:        maxPlayerCount,
		Analytics:             true,
		RecordTime:            true,
		ViewOpponentsSessions: true,
	}
}

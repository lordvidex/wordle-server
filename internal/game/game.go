package game

import (
	"time"

	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
)

// Event is an enum for all the events that Listeners need to subscribe to
// when listening to updates to the Game state
type Event string

const (
	EventPlayerJoined  Event = "PlayerJoined"
	EventPlayerLeft    Event = "PlayerLeft"
	EventGameStarted   Event = "GameStarted"
	EventGameEnded     Event = "GameEnded"
	EventPlayerGuessed Event = "PlayerGuessed"
)

const (
	// MaxDuration is the maximum duration a game can last
	MaxDuration = time.Hour
)

type Game struct {
	// ID is the app specific identifier for a game in the app
	ID uuid.UUID

	// InviteID is the id of the lobby used to play this game
	InviteID string

	// Word is the correct word that should be guessed
	Word words.Word

	// Sessions represent each player's game session
	Sessions []*Session

	// Settings represent the rules of the game as set by the room owner
	Settings Settings

	// StartTime is the time the game started,
	// when the value is nil, this means the game has not started
	StartTime *time.Time

	// EndTime is the time the game ended,
	// when the value is nil, this means the game has not ended
	EndTime *time.Time
}

// HasEnded Game if the Game.EndTime is set OR if the game has been active for an hour
// Ended games do not receive rewards after completed Sessions and penalties are applied
// to all sessions immediately after Game has ended.
// or if they have guessed the word correctly
func (g *Game) HasEnded() bool {
	if g.HasStarted() {
		return false
	}
	return g.EndTime != nil && g.EndTime.After(g.StartTime.Add(MaxDuration))
}

func (g *Game) HasStarted() bool {
	return g.StartTime != nil
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
	PlayerCount           int  `json:"player_count"`
	Analytics             bool `json:"has_analytics"`
	RecordTime            bool `json:"should_record_time"`
	ViewOpponentsSessions bool `json:"can_view_opponents_sessions"`
}

// NewSettings creates a new Settings struct
// with default values
func NewSettings(playerCount int) Settings {
	return Settings{
		WordLength:            5,
		Trials:                6,
		PlayerCount:           playerCount, 
		Analytics:             true,
		RecordTime:            true,
		ViewOpponentsSessions: true,
	}
}

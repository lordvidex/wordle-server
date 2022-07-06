package game

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
	"reflect"
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

type Game struct {
	ID uuid.UUID
	// Word is the correct word that should be guessed
	Word words.Word
	// PlayerSessions represent each player's game session
	PlayerSessions map[Player]*Session
	// Settings represent the rules of the game as set by the room owner
	Settings Settings
}

// HasEnded Game if all the players in the game have used up all their guesses
// or if they have guessed the word correctly
func (g Game) HasEnded() bool {
	for _, session := range g.PlayerSessions {
		if !session.Complete(g.Settings.Trials, g.Word) {
			return false
		}
	}
	return true
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

func (s *Session) Complete(maxTries int, correct words.Word) bool {
	if len(s.Guesses) == 0 {
		return false
	}
	lastGuess := s.Guesses[len(s.Guesses)-1]
	isEqual := reflect.DeepEqual(lastGuess.Letters.Keys(), correct.Letters.Keys())
	if len(s.Guesses) >= maxTries || isEqual {
		return true
	}
	return false
}

type Settings struct {
	WordLength            int  `json:"word_length"`
	Trials                int  `json:"trials"`
	PlayerCount           int  `json:"player_count"`
	Analytics             bool `json:"has_analytics"`
	RecordTime            bool `json:"should_record_time"`
	ViewOpponentsSessions bool `json:"should_view_opponents_sessions"`
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

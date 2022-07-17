package game

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
	"time"
)

const (
	RandomWordLength = 5
)

type CreateGameCommand struct {
	CaptainID   uuid.UUID
	Settings    Settings
	InviteID    string
	PlayerCount int
}

type CreateGameHandler interface {
	Handle(command CreateGameCommand) (*Game, error)
}

type createGameHandler struct {
	repo                Repository
	randomWordGenerator words.RandomHandler
}

func NewCreateGameHandler(repo Repository, randomWordGen words.RandomHandler) CreateGameHandler {
	return &createGameHandler{repo, randomWordGen}
}

func (h *createGameHandler) Handle(command CreateGameCommand) (*Game, error) {
	game := &Game{
		ID:          uuid.New(),
		Settings:    command.Settings,
		Word:        h.randomWordGenerator.GetRandomWord(RandomWordLength),
		InviteID:    command.InviteID,
		StartTime:   time.Now(),
		CreatorID:   command.CaptainID,
		PlayerCount: command.PlayerCount,
	}
	game, err := h.repo.Create(game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

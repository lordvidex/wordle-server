package game

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
)

type StartGameCommand struct {
	Players  []Player
	Settings Settings
}
type StartGameCommandHandler interface {
	Handle(command StartGameCommand) error
}

type startGameCommandHandler struct {
	repo     Repository
	wordGen  words.RandomHandler
	notifier NotificationService
}

func (c *startGameCommandHandler) Handle(command StartGameCommand) error {
	sessions := make(map[Player]*Session)
	for _, player := range command.Players {
		sessions[player] = NewSession(&player)
	}
	game := &Game{
		ID:             uuid.New(),
		PlayerSessions: sessions,
		Settings:       command.Settings,
		Word:           c.wordGen.GetRandomWord(command.Settings.WordLength),
	}
	game, err := c.repo.Create(game)
	if err != nil {
		return err
	}
	return c.notifier.UpdateGameState(EventGameStarted, game)

}

func NewStartGameCommandHandler(repo Repository, g words.RandomHandler, n NotificationService) StartGameCommandHandler {
	return &startGameCommandHandler{repo, g, n}
}

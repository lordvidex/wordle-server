package game

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
)

type StartGameCommand struct {
	ID uuid.UUID
}
type StartGameHandler interface {
	Handle(command StartGameCommand) error
}

type startGameCommandHandler struct {
	repo     Repository
	wordGen  words.RandomHandler
	notifier NotificationService
}

func (c *startGameCommandHandler) Handle(command StartGameCommand) error {
	game, err := c.repo.FindByID(command.ID.String(), &Player{})
	if err != nil {
		return err
	}
	err = c.repo.Start(game.ID.String())
	if err != nil {
		return err
	}
	return c.notifier.UpdateGameState(EventGameStarted, game)

}

func NewStartGameCommandHandler(repo Repository, g words.RandomHandler, n NotificationService) StartGameHandler {
	return &startGameCommandHandler{repo, g, n}
}

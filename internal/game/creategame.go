package game

import "github.com/google/uuid"

type CreateGameCommand struct {
	Player   *Player
	Settings Settings
}

type CreateGameHandler interface {
	Handle(command CreateGameCommand) (*Game, error)
}

type createGameHandler struct {
	repo              Repository
	inviteIdGenerator InviteIDGenerator
}

func NewCreateGameHandler(repo Repository, inviteIdGenerator InviteIDGenerator) CreateGameHandler {
	return &createGameHandler{repo, inviteIdGenerator}
}

func (h *createGameHandler) Handle(command CreateGameCommand) (*Game, error) {
	game := &Game{
		ID:       uuid.New(),
		InviteID: h.inviteIdGenerator.Generate(),
		Settings: command.Settings,
	}
	game, err := h.repo.Create(game)
	if err != nil {
		return nil, err
	}
	// add the player that created the game
	err = h.repo.Join(game.ID.String(), command.Player)
	if err != nil {
		return nil, err
	}
	return game, nil
}

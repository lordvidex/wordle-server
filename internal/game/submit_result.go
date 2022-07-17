package game

import (
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/words"
)

type SubmitResultCommand struct {
	GameID uuid.UUID
	PlayerID uuid.UUID
	PlayerName string
	Words []*words.Word
}

type SubmitResultHandler interface {
	Handle(command SubmitResultCommand) error
}

type submitResultHandler struct {
	repo Repository
	AwardSystem AwardSystem
}

func NewSubmitResultHandler(repo Repository, awardSystem AwardSystem) SubmitResultHandler {
	return &submitResultHandler{repo, awardSystem}
}

func (h *submitResultHandler) Handle(command SubmitResultCommand) error {
	return h.repo.UpdateGameResult(
		command.GameID,
		command.PlayerID,
		command.PlayerName, 
		h.AwardSystem.AwardPoints,
		command.Words,
	)
}
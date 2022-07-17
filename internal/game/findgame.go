package game

import (
	"errors"

	"github.com/google/uuid"
)

type FindGameQuery struct {
	ID     uuid.UUID
	Invite string
}

type FindGameByIDQueryHandler interface {
	Handle(query FindGameQuery) (*Game, error)
}

type findGameByIDQueryHandler struct {
	repo Repository
}

func (h *findGameByIDQueryHandler) Handle(query FindGameQuery) (*Game, error) {
	if query.ID == uuid.Nil {
		return nil, errors.New("game id is required")
	}
	return h.repo.FindByID(query.ID)
}

func NewFindGameByIDQueryHandler(repo Repository) FindGameByIDQueryHandler {
	return &findGameByIDQueryHandler{repo}
}
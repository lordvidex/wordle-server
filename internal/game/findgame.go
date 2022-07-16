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
	return h.repo.FindByID(query.ID.String())
}

func NewFindGameByIDQueryHandler(repo Repository) FindGameByIDQueryHandler {
	return &findGameByIDQueryHandler{repo}
}

type FindAllGamesQueryHandler interface {
	Handle() ([]Game, error)
}

type findAllGamesQueryHandler struct {
	repo Repository
}

func (h *findAllGamesQueryHandler) Handle() ([]Game, error) {
	return h.repo.FindAll()
}

func NewFindAllGamesQueryHandler(repo Repository) FindAllGamesQueryHandler {
	return &findAllGamesQueryHandler{repo}
}

type FindByInviteIDQueryHandler interface {
	Handle(query FindGameQuery) ([]Game, error)
}

type findByInviteIDQueryHandler struct {
	repo Repository
}

func (h *findByInviteIDQueryHandler) Handle(query FindGameQuery) ([]Game, error) {
	if query.Invite == "" {
		return nil, errors.New("invite is required")
	}
	return h.repo.FindByInviteID(query.Invite)
}

func NewFindByInviteIDQueryHandler(repo Repository) FindByInviteIDQueryHandler {
	return &findByInviteIDQueryHandler{repo}
}

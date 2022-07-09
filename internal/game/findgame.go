package game

import "github.com/google/uuid"

type FindGameQuery struct {
	ID uuid.UUID
}

type FindGameQueryHandler interface {
	Handle(query FindGameQuery) (*Game, error)
}

type findGameQueryHandler struct {
	repo Repository
}

func (h *findGameQueryHandler) Handle(query FindGameQuery) (*Game, error) {
	return h.repo.Find(query.ID.String())
}

func NewFindGameQueryHandler(repo Repository) FindGameQueryHandler {
	return &findGameQueryHandler{repo}
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

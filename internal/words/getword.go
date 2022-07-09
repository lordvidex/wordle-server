package words

import "github.com/google/uuid"

type GetWordQuery struct {
	ID uuid.UUID
}

type GetWordHandler interface {
	Handle(query GetWordQuery) (*Word, error)
}

type getWordQueryHandler struct {
	repo Repository
}

func (h *getWordQueryHandler) Handle(query GetWordQuery) (*Word, error) {
	return h.repo.Find(query.ID)
}

func NewGetWordHandler(repo Repository) GetWordHandler {
	return &getWordQueryHandler{repo}
}

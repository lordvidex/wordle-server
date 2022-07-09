package words

type AddWordCommand struct {
	Word Word
}
type AddStringWordCommand struct {
	Word string
}

type AddWordHandler interface {
	Handle(command AddWordCommand) error
	HandleString(command AddStringWordCommand) error
}

type addWordHandler struct {
	repo Repository
}

func NewAddWordHandler(repo Repository) AddWordHandler {
	return &addWordHandler{repo}
}
func (h *addWordHandler) Handle(command AddWordCommand) error {
	return h.repo.Add(command.Word)
}
func (h *addWordHandler) HandleString(command AddStringWordCommand) error {
	word := NewFromString(command.Word)
	wCommand := AddWordCommand{Word: *word}
	return h.Handle(wCommand)
}

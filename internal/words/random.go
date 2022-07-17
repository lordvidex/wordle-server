package words

// RandomHandler is an interface for getting a random word
// given the length
type RandomHandler interface {
	GetRandomWord(length int) Word
}

// StringGenerator is an interface for generating random strings
// of specific lengths
type StringGenerator interface {
	Generate(length int) string
}

type randomWordHandler struct {
	generator StringGenerator
}

func (h *randomWordHandler) GetRandomWord(length int) Word {
	return New(h.generator.Generate(length))
}

func NewRandomHandler(generator StringGenerator) RandomHandler {
	return &randomWordHandler{generator}
}

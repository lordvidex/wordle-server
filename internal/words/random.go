package words

import "github.com/google/uuid"

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
	// generate a random string
	word := h.generator.Generate(length)

	// map all the letters to correct
	letters := make(Letters, len([]rune(word)))
	for i, letter := range word {
		letters[i] = &Letter{
			Rune:   letter,
			Status: Correct,
		}
	}

	// return word
	return Word{
		ID:      uuid.New(),
		Letters: letters,
	}
}

func NewRandomHandler(generator StringGenerator) RandomHandler {
	return &randomWordHandler{generator}
}

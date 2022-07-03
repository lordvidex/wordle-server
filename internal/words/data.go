package words

import (
	_ "embed"
	"github.com/lordvidex/wordle-wf/internal/common/werr"
	"log"
	"strconv"
	"strings"
)

var (
	//go:embed resources/five_letter_words.txt
	fileContent string
)

type repository struct {
	words   map[int]string
	reverse map[string]int
	count   int
}

func (r *repository) AddWord(s string) (int, error) {
	r.reverse[s] = r.count
	r.words[r.count] = s
	r.count++
	return r.count - 1, nil
}

func (r *repository) GetWordByID(i int) (*Word, error) {
	if word, exists := r.words[i]; !exists {
		return nil, werr.NotFound("No word with id " + strconv.Itoa(i))
	} else {
		return &Word{i, word}, nil
	}
}

func (r *repository) GetWordByValue(s string) (*Word, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) SaveWord(word Word) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) GetWordCount() int {
	return r.count
}

func (r *repository) loadFiveLetterWords() {
	// read the file's content into our array
	words := strings.Split(fileContent, "\n")
	for _, word := range words {
		_, err := r.AddWord(word)
		if err != nil {
			log.Fatal("error adding word ", word, ": ", err)
		}
	}
}

func NewRepository() Repository {
	repo := &repository{
		words:   make(map[int]string),
		reverse: make(map[string]int),
		count:   0,
	}
	repo.loadFiveLetterWords()
	return repo
}

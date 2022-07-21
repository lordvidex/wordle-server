package adapters

import (
	// embed here helps to embed the words in the resources folder
	// into fileContent
	_ "embed"
	"github.com/lordvidex/wordle-wf/internal/words"
	"math/rand"
	"strings"
	"time"
)

var (
	//go:embed resources/five_letter_words.txt
	fileContent string
)

type localWordGenerator struct {
	Words []string
}

func (l *localWordGenerator) Generate(length int) string {
	if length != 5 {
		panic("only 5 letter words are supported")
	}
	rand.Seed(time.Now().UnixNano())
	return l.Words[rand.Intn(len(l.Words))]
}

func (l *localWordGenerator) loadWords() {
	// read the file's content into our array
	l.Words = strings.Split(fileContent, "\n")
}

func NewLocalStringGenerator() words.StringGenerator {
	l := &localWordGenerator{}
	l.loadWords()
	return l
}

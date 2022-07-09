package words

import (
	"time"

	"github.com/google/uuid"
)

// LetterStatus is an enum type for the Status of a letter in a word guess
// -1 stands for Incorrect
// 0 stands for a letter that is in the wrong position
// 1 stands for a correct letter in the right position
type LetterStatus int

const (
	Incorrect LetterStatus = iota - 1
	Exists
	Correct
)

// Letter is a struct that represents a letter in a word guess and it's Status
type Letter struct {
	Rune   rune         `json:"rune"`
	Status LetterStatus `json:"status"`
}

// Letters represent a struct of Letter which provides functions to conveniently
// interact with a word as a whole
type Letters []*Letter

// Keys returns the keys of the Letters map
func (l Letters) Keys() []rune {
	runes := make([]rune, len(l))
	i := 0
	for _, v := range l {
		runes[i] = v.Rune
		i++
	}
	return runes
}

// Values returns the values of the Letters map
func (l Letters) Values() []LetterStatus {
	values := make([]LetterStatus, len(l))
	i := 0
	for _, v := range l {
		values[i] = v.Status
		i++
	}
	return values
}

// Word contains a map of letters to their Status
// and the time this word was played
// for example the word 'WEIRD' would have the following
// Letters mapping
//
// W -> Incorrect
// E -> Correct
// I -> Incorrect
// R -> Exists
// D -> Incorrect
//
type Word struct {
	ID         uuid.UUID
	Letters    Letters
	TimePlayed time.Time
}

func NewFromString(word string) *Word {
	runes := make(Letters, len([]rune(word)))
	for i, char := range word {
		runes[i] = &Letter{
			Rune:   char,
			Status: Incorrect,
		}
	}
	return &Word{
		ID:      uuid.New(),
		Letters: runes,
	}
}

// SetLetterStatus compares the word to the correct word
// and sets the Status of each letter of *Word accordingly
// Space Complexity: O(n)
// Time Complexity: O(n)
func (w *Word) SetLetterStatus(correctWord *Word) {
	// change all LetterStatuses to Incorrect
	for key := range w.Letters {
		w.Letters[key].Status = Incorrect
	}

	// check if the lengths match
	if len(w.Letters) != len(correctWord.Letters) {
		return
	}

	// make a dict of the correct letters
	dict := make(map[rune]int)
	for _, v := range correctWord.Letters {
		dict[v.Rune] += 1
	}

	// first parse the correct letters
	for i, v := range w.Letters {
		if v.Rune == correctWord.Letters[i].Rune {
			v.Status = Correct
			dict[v.Rune] -= 1
		}
	}

	// parse the letters that have wrong positions
	for _, value := range w.Letters {
		if value.Status == Correct {
			continue
		}
		if cnt, ok := dict[value.Rune]; ok && cnt > 0 {
			value.Status = Exists
			dict[value.Rune] -= 1
		}
	}
}

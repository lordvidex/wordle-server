package words

import (
	"database/sql"
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
	Word     string
	PlayedAt sql.NullTime
}

func New(word string) Word {
	return Word{word, sql.NullTime{}}
}

func (w Word) Runes() []rune {
	return []rune(w.Word)
}

// CompareTo compares the word to the correct word
// and returns LetterStatus of each letter of Word accordingly
// Space Complexity: O(n)
// Time Complexity: O(n)
func (w Word) CompareTo(correctWord Word) []LetterStatus {
	correctRunes := correctWord.Runes()
	instanceRunes := w.Runes()

	wordStatus := make([]LetterStatus, len(instanceRunes))
	for key := range instanceRunes {
		wordStatus[key] = Incorrect
	}

	// check if the lengths match
	if len(instanceRunes) != len(correctRunes) {
		return wordStatus
	}

	// make a dict of the correct letters
	dict := make(map[rune]int)
	for _, v := range correctRunes {
		dict[v] += 1
	}

	// first parse the correct letters
	for i, v := range instanceRunes {
		if v == correctRunes[i] {
			wordStatus[i] = Correct
			dict[v] -= 1
		}
	}

	// parse the letters that have wrong positions
	for i, value := range instanceRunes {
		if wordStatus[i] == Correct {
			continue
		}
		if cnt, ok := dict[value]; ok && cnt > 0 {
			wordStatus[i] = Exists

			dict[value] -= 1
		}
	}
	return wordStatus
}

func (w Word) String() string {
	return w.Word
}

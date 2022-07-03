package words

import (
	"testing"
)

type void struct{}

var member void

func TestDataIsReadFromJSON(t *testing.T) {
	if len(words) == 0 {
		t.Fail()
	}
}

func TestGetRandomWord(t *testing.T) {
	t.Run("should return a random word 6 times", func(t *testing.T) {
		wordSet := make(map[string]void)
		for i := 0; i < 6; i++ {
			wordSet[GetRandomWord()] = member
		}
		if len(wordSet) != 6 {
			t.Fail()
		}
	})
}

package game

import (
	"github.com/lordvidex/wordle-wf/internal/words"
	"testing"
)

func TestSession_Complete(t *testing.T) {
	corr := words.New("HELLO")
	testCases := []struct {
		Guesses     []words.Word
		Tries       int
		Expected    bool
		description string
	}{
		{[]words.Word{}, 6, false, "no guesses"},
		{[]words.Word{
			words.New("WORDS"),
			words.New("HELLO"),
		}, 6, true, "correct after two guesses"},
		{[]words.Word{
			words.New("WORDS"),
			words.New("HALLO"),
		}, 2, true, "max tries"},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			s := Session{
				Guesses: tt.Guesses,
			}
			if got := s.Complete(tt.Tries, *corr); got != tt.Expected {
				t.Errorf("Session.Complete() = %v, want %v", got, tt.Expected)
			}
		})
	}
}

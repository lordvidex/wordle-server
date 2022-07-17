package app

import (
	"github.com/lordvidex/wordle-wf/internal/words"
	"testing"
)

func TestSession_IsWon(t *testing.T) {
	tests := []struct {
		name    string
		session *Session
		want    bool
	}{
		{"no sessions", &Session{}, false},
		{"no words", &Session{playedWords: []words.Word{}}, false},
		{"last word is correct", &Session{playedWords: []words.Word{
			words.Word("a"),
		},
			correctWord: "a",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.session.IsWon(); got != tt.want {
				t.Errorf("IsWon() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mock struct{}

func (mock) GetRandomWord(length int) words.Word {
	return words.New("CORRECT")
}
func TestNewSessionGeneratesWord(t *testing.T) {
	// given
	tries := 5
	s := NewSession(tries, mock{})
	if s.correctWord == "" {
		t.Error("NewSession() did not generate a word")
	}
	if s.maxTries != tries {
		t.Errorf("NewSession() did not set maxTries to %d", tries)
	}
}

func TestSession_HasEnded(t *testing.T) {
	type fields struct {
		maxTries    int
		playedWords []words.Word
		correctWord words.Word
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"new session",
			fields{maxTries: 10, playedWords: []words.Word{}, correctWord: words.New("CORRECT")},
			false,
		},
		{"finished session",
			fields{maxTries: 4, playedWords: []words.Word{
				words.New("FIRST"),
				words.New("SECOND"),
				words.New("THIRD"),
				words.New("FOURTH"),
			}, correctWord: "CORRECT"}, true},
		{"correct word", fields{maxTries: 4, playedWords: []words.Word{
			words.New("FIRST"),
			words.New("CORRECT"),
		}, correctWord: words.New("CORRECT"),
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				maxTries:    tt.fields.maxTries,
				playedWords: tt.fields.playedWords,
				correctWord: tt.fields.correctWord,
			}
			if got := s.HasEnded(); got != tt.want {
				t.Errorf("HasEnded() = %v, want %v", got, tt.want)
			}
		})
	}
}

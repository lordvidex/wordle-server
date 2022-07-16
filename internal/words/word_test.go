package words

import (
	"reflect"
	"testing"
)

func TestLetterStateEnums(t *testing.T) {
	testCases := []struct {
		letterStatus LetterStatus
		expected     int
		errorMessage string
	}{
		{Correct, 1, "Correct should be 1"},
		{Incorrect, -1, "Incorrect should be -1"},
		{Exists, 0, "Exists should be 0"},
	}

	for _, tt := range testCases {
		t.Run(tt.errorMessage, func(t *testing.T) {
			if int(tt.letterStatus) != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, tt.letterStatus)
			}
		})
	}
}

func TestWord_SetLetterStatus(t *testing.T) {
	testCases := []struct {
		word        string
		correctWord string
		expected    []LetterStatus
		desc        string
	}{
		{"WEIRD", "WORLD", []LetterStatus{Correct, Incorrect, Incorrect, Exists, Correct}, "contains WRD"},
		{"SAVED", "WORLD", []LetterStatus{Incorrect, Incorrect, Incorrect, Incorrect, Correct}, "contains just D"},
		{"SEIZE", "WORLD", []LetterStatus{Incorrect, Incorrect, Incorrect, Incorrect, Incorrect}, "contains nothing"},
		{"SEGMENT", "WORLD", []LetterStatus{Incorrect, Incorrect, Incorrect, Incorrect, Incorrect, Incorrect, Incorrect}, "longer than word to be guessed"},
		{"SEX", "WORLD", []LetterStatus{Incorrect, Incorrect, Incorrect}, "shorter than word to be guessed"},
		{"LOROC", "WORLD", []LetterStatus{Exists, Correct, Correct, Incorrect, Incorrect}, "One correct 'O' and One wrong 'O'"},
		{"ALELE", "EVENT", []LetterStatus{Incorrect, Incorrect, Correct, Incorrect, Exists}, "One correct E and One wrong E"},
		{"EVENT", "EVENT", []LetterStatus{Correct, Correct, Correct, Correct, Correct}, "Same word"},
		{"RITES", "SITES", []LetterStatus{Incorrect, Correct, Correct, Correct, Correct}, "Wrong letter first that exists later"},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			// given
			word := NewFromString(tt.word)
			correctWord := NewFromString(tt.correctWord)
			//when
			word.SetLetterStatus(correctWord)
			// then
			if !reflect.DeepEqual(word.Letters.Values(), tt.expected) {
				t.Errorf("Expected %v [%v], got %v [%v]",
					tt.expected,
					tt.correctWord,
					word.Letters.Values(),
					tt.word)
			}
		})
	}
}

func TestLetters_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		l       Letters
		args    args
		wantErr bool
	}{
		{name: "bytes test", l: Letters{}, args: args{src: []byte(`{"W": -1, "E": -1, "I": -1, "R": -1, "D": -1}`)}, wantErr: false},
		{name: "string test", l: Letters{}, args: args{src: `{"W": -1, "E": -1, "I": -1, "R": -1, "D": -1}`}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if string(tt.l.Keys()) != "WEIRD" {
				t.Errorf("Expected %s, got %v", "WEIRD", tt.l.Keys())
			}
			if len(tt.l.Values()) != 5 {
				t.Errorf("values not scanned properly, expected 5, got %v", len(tt.l.Values()))
			}
		})
	}
}

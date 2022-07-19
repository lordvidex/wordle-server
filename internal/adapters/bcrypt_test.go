package adapters

import (
	"testing"

	"github.com/lordvidex/wordle-wf/internal/auth"
)

func computeHash(h auth.PasswordHelper, value string, t *testing.T) string {
	res, err := h.Hash(value)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	return res
}

func TestBcryptHandler_Validate(t *testing.T) {
	// given
	h := NewBcryptHelper()

	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name: "normal password",
			password: "password",
			hash: computeHash(h, "password", t),
			expected: true,
		},
		{
			name: "short password",
			password: "pass",
			hash: computeHash(h, "pass", t),
			expected: true,
		},
		{
			name: "long password",
			password: "asdgburgbewirthgewin4",
			hash: computeHash(h, "asdgburgbewirthgewin4", t),
			expected: true,
		}, 
		{
			name: "wrong password",
			password: "wrong",
			hash: computeHash(h, "password", t),
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			res := h.Validate(tt.password, tt.hash)

			// then
			if res != tt.expected {
				t.Errorf("Expected true for validating %s text and %s hash, got false", tt.password, tt.hash)
			}
		})
	}
}

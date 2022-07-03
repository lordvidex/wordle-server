package words

import (
	"encoding/json"
	"io"
)

type Word struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

func (word *Word) WriteJSON(w io.Writer) {
	err := json.NewEncoder(w).Encode(word)
	if err != nil {
		panic(err)
	}
}

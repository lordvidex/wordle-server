package werr

import (
	"encoding/json"
	"io"
)

type WordErr struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *WordErr) Error() string {
	return e.Message
}

func (e *WordErr) WriteJSON(w io.Writer) {
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		panic(err)
	}
}

// NotFound is a convenient function for returning a 404 error with a message
func NotFound(message string) *WordErr {
	if message == "" {
		message = "Not Found"
	}
	return &WordErr{404, message}
}

// BadRequest is a convenient function for returning a 400 error with a message
func BadRequest(message string) *WordErr {
	return &WordErr{400, message}
}

// InternalServerError is a convenient function for returning a 500 error with a message
// can be converted to werr.Error() if a thrown error is the cause
func InternalServerError(message string) *WordErr {
	return &WordErr{500, message}
}

// NewWordErr receives
// code - as http status code
// and
// message - custom status message
func NewWordErr(code int, message string) *WordErr {
	return &WordErr{code, message}
}

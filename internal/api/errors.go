package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) WriteJSON(w io.Writer) {
	_ = json.NewEncoder(w).Encode(e)
}

// NotFound is a convenient function for returning a 404 error with a message
func NotFound(message string) Error {
	if message == "" {
		message = "Not Found"
	}
	return Error{http.StatusNotFound, message}
}

// BadRequest is a convenient function for returning a 400 error with a message
func BadRequest(message string) Error {
	return Error{http.StatusBadRequest, message}
}

// Unauthorized is a convenient function for returning a 401 error with a message
func Unauthorized(message string) Error {
	return Error{http.StatusUnauthorized, message}
}

func Forbidden(message string) Error {
	return Error{http.StatusForbidden, message}
}

// InternalServerError is a convenient function for returning a 500 error with a message
// can be converted to werr.Error() if a thrown error is the cause
func InternalServerError(message string) Error {
	return Error{http.StatusInternalServerError, message}
}

// NewError receives
// code - as http status code
// and
// message - custom status message
func NewError(code int, message string) Error {
	return Error{code, message}
}

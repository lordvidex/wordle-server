package middleware

import (
	"fmt"
	"github.com/lordvidex/wordle-wf/internal/common/werr"
	"net/http"
)

type ResponseStatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (w ResponseStatusRecorder) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

// JSONContent is a middleware that sets the Content-Type header to application/json
func JSONContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := ResponseStatusRecorder{ResponseWriter: w, Status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		fmt.Printf("RESPONSE [%s] %d %s\n", r.Method, recorder.Status, r.URL.Path)
	})
}

func HandleError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				wordErr, ok := err.(*werr.WordErr)
				if !ok {
					wordErr = werr.InternalServerError(err.(error).Error())
				}
				w.WriteHeader(wordErr.StatusCode)
				wordErr.WriteJSON(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

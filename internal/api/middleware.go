package api

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	"github.com/lordvidex/wordle-wf/internal/auth"
)

type AuthContextKey string

const (
	DecodedUserKey AuthContextKey = "user_payload"
)
const (
	AuthHeader            = "Authorization"
	AuthHeaderValuePrefix = "Bearer"
)

var (
	ErrBadToken      = errors.New("invalid or expired token")
	ErrNoToken       = errors.New("no token provided")
	ErrBadAuthHeader = errors.New("authorization header is badly formatted")
)

type responseStatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (w *responseStatusRecorder) WriteHeader(status int) {
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
		recorder := responseStatusRecorder{ResponseWriter: w, Status: http.StatusOK}
		next.ServeHTTP(&recorder, r)
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": recorder.Status,
		}).Info("RESPONSE")
	})
}

func HandleError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var wordErr Error
				switch err.(type) {
				case error:
					wordErr = InternalServerError(err.(error).Error())
				case string:
					wordErr = InternalServerError(err.(string))
				default:
					wordErr = InternalServerError("unknown error occurred")
				}
				w.WriteHeader(wordErr.StatusCode)
				wordErr.WriteJSON(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(tokenDecoder auth.GetUserByTokenQueryHandler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(AuthHeader)
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				Unauthorized(ErrNoToken.Error()).WriteJSON(w)
				return
			}

			fields := strings.Fields(token)
			if len(fields) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				BadRequest(ErrBadAuthHeader.Error()).WriteJSON(w)
				return
			}

			if !strings.EqualFold(fields[0], AuthHeaderValuePrefix) {
				w.WriteHeader(http.StatusBadRequest)
				BadRequest(ErrBadAuthHeader.Error()).WriteJSON(w)
				return
			}

			player, err := tokenDecoder.Handle(auth.Token(fields[1]))
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				Unauthorized(ErrBadToken.Error()).WriteJSON(w)
				return
			}
			ctx := context.WithValue(r.Context(), DecodedUserKey, player)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

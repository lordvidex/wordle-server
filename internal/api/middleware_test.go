package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/auth"
	"github.com/lordvidex/wordle-wf/internal/game"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	tokenAdapter := auth.NewMockTokenHelper(ctrl)

	testPlayer := &game.Player{
		ID:    uuid.New(),
		Name:  "mathew",
		Email: "test@gmail.com",
	}
	tests := []struct {
		name    string
		prepare func(r *http.Request)
		want    *game.Player
		wantErr bool
	}{
		{"no token", func(r *http.Request) {}, nil, true},
		{"invalid token",
			func(r *http.Request) {
				r.Header.Set("Authorization", "Bearer invalid")
				tokenAdapter.EXPECT().Decode(auth.Token("invalid"), gomock.Any()).Return(auth.ErrInvalidToken)
			}, nil, true},
		{"valid token", func(r *http.Request) {
			r.Header.Set("Authorization", "Bearer valid")
			tokenAdapter.EXPECT().Decode(auth.Token("valid"), gomock.Any()).Return(nil)
		}, testPlayer, false},
		{"no bearer", func(r *http.Request) {
			r.Header.Set("Authorization", "valid")
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			r := httptest.NewRequest("GET", "/", nil)
			tt.prepare(r)
			recorder := httptest.NewRecorder()
			middleware := AuthMiddleware(auth.NewUserTokenDecoder(tokenAdapter))
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("Test handler receives: ", r.Context().Value(DecodedUserKey))
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("Authenticated"))
			})
			handler := middleware(testHandler)

			// when
			handler.ServeHTTP(recorder, r)

			// then
			if tt.wantErr {
				if recorder.Code == http.StatusOK {
					t.Errorf("expected status code to be 400...500, got %d", recorder.Code)
				}
				fmt.Println("TEST Passed with payload", recorder.Body.String())
			} else {
				val := r.Context().Value(DecodedUserKey)
				if val == nil {
					t.Errorf("expected user to be set in context")
				}
				_, ok := val.(*game.Player)
				if !ok {
					t.Errorf("user not decoded as game.Player")
				}
			}
		})
	}
}

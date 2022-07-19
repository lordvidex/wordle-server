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
			}, nil, true},
		{"valid token", func(r *http.Request) {
			r.Header.Set("Authorization", "Bearer valid")
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
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("Authenticated"))
			})
			handler := middleware(testHandler)

			// when
			handler.ServeHTTP(recorder, r)

			// then
			if tt.wantErr {
				if recorder.Code != http.StatusUnauthorized {
					t.Errorf("expected status code 401, got %d", recorder.Code)
				}
				fmt.Println("TEST Passed with payload", recorder.Body.String())
			} else {
				val := r.Context().Value(DecodedUserKey)
				if val == nil {
					t.Errorf("expected user to be set in context")
				}
				valPlayer, ok := val.(*game.Player)
				if !ok {
					t.Errorf("user not decoded as game.Player")
				}
				if valPlayer.ID != testPlayer.ID {
					t.Errorf("expected user to have ID %v, got %v", testPlayer.ID, valPlayer.ID)
				}
				if valPlayer.Name != testPlayer.Name {
					t.Errorf("expected user to have name %v, got %v", testPlayer.Name, valPlayer.Name)
				}
				if valPlayer.Email != testPlayer.Email {
					t.Errorf("expected user to have email %v, got %v", testPlayer.Email, valPlayer.Email)
				}
			}
		})
	}
}

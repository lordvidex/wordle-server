package websockets

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/api"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/pkg/errors"
)

var (
	testUUID = uuid.New()
)

func addPlayerToReq(r *http.Request) {
	player := &game.Player{
		ID:        testUUID,
		Name:      "James",
		Points:    20,
		Email:     "test@gmail.com",
		IsDeleted: false,
	}
	*r = *r.WithContext(context.WithValue(r.Context(), api.DecodedUserKey, player))
}
func TestJoinLobby(t *testing.T) {
	testCases := []struct {
		name               string
		rooms              map[string]*Room
		prepare            func(r *http.Request, t *testing.T)
		expectError        bool
		expectedStatusCode int
		error              error
	}{
		{"unauthenticated user",
			nil,
			func(*http.Request, *testing.T) {},
			true,
			http.StatusUnauthorized,
			ErrPlayerNotInRequest,
		},
		{"no room id",
			nil,
			func(r *http.Request, _ *testing.T) {
				addPlayerToReq(r)
			},
			true,
			http.StatusBadRequest,
			ErrRoomNotInRequest,
		},
		{"room id not found",
			func() map[string]*Room {
				rr := make(map[string]*Room)
				rr["test"] = NewRoom("test", game.Settings{})
				return rr
			}(),
			func(r *http.Request, t *testing.T) {
				addPlayerToReq(r)
				parsedURL, err := url.Parse("/live?room_id=123")
				if err != nil {
					t.Fatal(err)
				}
				r.URL = parsedURL
			},
			true,
			http.StatusNotFound,
			fmt.Errorf("room with id: 123 does not exist"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/live", nil)
			r.URL.Scheme = "ws"

			tt.prepare(r, t)
			g := &GameSocket{
				rooms: tt.rooms,
			}
			defer g.Close()

			g.ServeHTTP(w, r)
			// s.ServeHTTP(w, r)
			if tt.expectError {
				if w.Code != tt.expectedStatusCode {
					t.Errorf("expected status code %v, got %v", tt.expectedStatusCode, w.Code)
				}
				var x map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &x)
				if err != nil {
					t.Fatal(errors.Wrap(err, "failed to unmarshal body response"))
				}
				if x["message"] != tt.error.Error() {
					t.Errorf("expected error message %v, got %v", tt.error.Error(), x["message"])
				}
			}
		})
	}
}

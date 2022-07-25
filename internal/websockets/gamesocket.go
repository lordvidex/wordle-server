package websockets

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/lordvidex/wordle-wf/internal/api"
	"github.com/lordvidex/wordle-wf/internal/game"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var (
	ErrRoomNotFound       = errors.New("game room not found")
	ErrPlayerNotInRequest = errors.New("player not found in request")
	ErrRoomNotInRequest   = errors.New("room id not in request")
	ErrCannotJoinLobby    = errors.New("cannot join lobby")
)

// others
var (
	queryRoomID     = "room_id"
	queryPlayerName = "player_name"
)

type GameSocket struct {
	rooms map[string]*Room
	Fgh   game.FindGameByIDQueryHandler
}

func (g *GameSocket) CreateLobby(settings *game.Settings, id string) (string, error) {
	g.rooms[id] = NewRoom(id, *settings)
	return id, nil
}

func (g *GameSocket) JoinLobby(lobbyID string, player *game.Player) {

}

func NewGameSocket(fgh game.FindGameByIDQueryHandler) *GameSocket {
	sock := &GameSocket{make(map[string]*Room), fgh}
	return sock
}

func (g *GameSocket) Close() error {
	var err error
	for id, room := range g.rooms {
		err = room.Close()
		delete(g.rooms, id)
	}
	return err
}

// ServeHTTP handles websocket requests for GameSocket - JoinGameEvent
// and connects the user to a game session if correct id is produced
func (g *GameSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	player, ok := ctx.Value(api.DecodedUserKey).(*game.Player)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		api.Unauthorized(ErrPlayerNotInRequest.Error()).WriteJSON(w)
		return
	}

	idList, ok := query[queryRoomID]
	if !ok || len(idList) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		api.BadRequest(ErrRoomNotInRequest.Error()).WriteJSON(w)
		return
	}

	id := idList[0] // get the first id
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		api.BadRequest(ErrRoomNotInRequest.Error()).WriteJSON(w)
		return
	}

	// check if the room exists
	room, exists := g.rooms[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		api.NotFound(fmt.Errorf("room with id: %v does not exist", id).Error()).WriteJSON(w)
		return
	}

	// check if the room has an active game
	shouldJoinGame := true
	if room.hasActiveGame {
		shouldJoinGame = false
		for key := range room.players {
			if key.playerID == player.ID.String() {
				shouldJoinGame = true
			}
		}
	}

	// set first player as owner
	if len(room.players) < 1 {
		room.owner = player.ID.String()
	}

	// upgrade the connection to a websocket
	if shouldJoinGame {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("error upgrading to websockets", err)
			return
		}
		room.join <- NewClient(room, conn, player.ID.String())
	} else {
		// throw a forbidden exception
		w.WriteHeader(http.StatusForbidden)
		api.Forbidden(ErrCannotJoinLobby.Error()).WriteJSON(w)
	}
}

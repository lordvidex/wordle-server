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
	ErrRoomNotFound = errors.New("game room not found")
)

// others
var (
	queryRoomID = "id"
)

type GameSocket struct {
	rooms map[string]*Room
	Fgh   game.FindGameByIDQueryHandler
}

func (g *GameSocket) CreateLobby(settings *game.Settings, id string) (string, error) {
	g.rooms[id] = NewRoom(id, *settings)
	return id, nil
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
	player := ctx.Value(api.DecodedUserKey).(*game.Player)

	idList, ok := r.URL.Query()[queryRoomID]
	if !ok || len(idList) < 1 {
		fmt.Println("error getting room id")
		return
	}

	id := idList[0] // get the first id
	room := g.rooms[id]
	canJoinGame := true
	// check if the room exists
	if room == nil {
		api.BadRequest(fmt.Errorf("room with id: %v does not exist", id).Error()).WriteJSON(w)
		return
	}
	// check if the room has an active game
	if room.hasActiveGame {
		canJoinGame = false
		for key := range room.players {
			if key.playerID == player.ID.String() {
				canJoinGame = true
			}
		}
	}
	// set first player as owner
	if len(room.players) < 1 {
		room.owner = player.ID.String()
	}

	// upgrade the connection to a websocket
	if canJoinGame {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("error upgrading to websockets", err)
			return
		}
		room.join <- NewClient(room, conn, player.ID.String())
	} else {
		// throw a forbidden exception
		api.Forbidden(fmt.Errorf("cannot join lobby").Error()).WriteJSON(w)
	}
}

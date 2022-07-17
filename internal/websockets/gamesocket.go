package websockets

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lordvidex/wordle-wf/internal/game"
)

// upgrader
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// errors
var (
	ErrRoomNotFound = errors.New("game room not found")
)

// others
var (
	queryGameID = "id"
)

type GameSocket struct {
	rooms map[string]*Room
	Fgh   game.FindGameByIDQueryHandler
}

func (g *GameSocket) CreateLobby(settings *game.Settings, id string) (string, error) {
	//TODO implement me
	panic("implement me")
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
	idList, ok := r.URL.Query()[queryGameID]
	if !ok || len(idList) < 1 {
		fmt.Println("error getting game id")
		return
	}

	id := idList[0]              // get the first id
	_uuid, err := uuid.Parse(id) // convert to UUID
	if err != nil {
		fmt.Println("error parsing game id")
		return
	}

	// upgrade the connection to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error upgrading to websockets", err)
		return
	}
	room, ok := g.rooms[id]
	if !ok {
		// check if such a gm exists
		gm, err := g.Fgh.Handle(game.FindGameQuery{ID: _uuid})
		if err != nil {
			fmt.Println("error finding game", err)
			return
		}
		// create new room
		g.rooms[id] = NewRoom(id, gm.Settings)
	}
	room.join <- NewClient(room, conn)
}

package websockets

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lordvidex/wordle-wf/internal/game"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	queryGameID = "id"
)

var (
	ErrRoomNotFound = errors.New("game room not found")
)

type GameSocket struct {
	rooms map[string]*Room
	Fgh   game.FindGameQueryHandler
}

func (g *GameSocket) Close() error {
	var err error
	for id, room := range g.rooms {
		err = room.Close()
		delete(g.rooms, id)
	}
	return err
}

func sessionSlice(sessionMap map[game.Player]*game.Session) []*game.Session {
	s := make([]*game.Session, len(sessionMap))
	i := 0
	for _, sess := range sessionMap {
		s[i] = sess
		i++
	}
	return s
}

func (g *GameSocket) UpdateGameState(ev game.Event, gm *game.Game) error {
	// check if room exists
	if room, ok := g.rooms[gm.ID.String()]; ok {
		// create the sessions slice from the map
		msg := WSGameDto{
			ID:       &gm.ID,
			Settings: &gm.Settings,
			Sessions: sessionSlice(gm.PlayerSessions),
		}
		payload := WSPayload{
			Event: ev,
			Data:  msg,
		}
		room.broadcast <- payload
		return nil
	} else {
		return ErrRoomNotFound
	}
}

// ServeHTTP handles websocket requests for GameSocket
// and connects the user to a game session if correct id is produced
func (g *GameSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idList, ok := r.URL.Query()["id"]
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
		// check if such a game exists
		_, err := g.Fgh.Handle(game.FindGameQuery{ID: _uuid})
		if err != nil {
			fmt.Println("error finding game", err)
			return
		}
		// create new room
		g.rooms[id] = NewRoom(id)
	}
	room.join <- NewClient(room, conn)
}

func NewGameSocket(fgh game.FindGameQueryHandler) *GameSocket {
	sock := &GameSocket{make(map[string]*Room), fgh}
	return sock
}

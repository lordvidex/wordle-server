package websockets

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lordvidex/wordle-wf/internal/game"
	"net/http"
	"strings"
	"time"
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

// WSGameDto is the game data payload sent to websockets clients
type WSGameDto struct {
	ID       *uuid.UUID      `json:"id"`
	Settings *game.Settings  `json:"settings"`
	Sessions []*game.Session `json:"sessions"`
}

// GameRoom is a room of connected and currently playing players
type GameRoom struct {
	ID      string
	players []*websocket.Conn
}

func (g *GameRoom) removePlayerAt(i int) {
	g.players[i] = g.players[len(g.players)-1]
	g.players[len(g.players)-1] = nil
	g.players = g.players[:len(g.players)-1]
}

func (g *GameRoom) Close() error {
	var err error
	for _, player := range g.players {
		err = player.Close()
	}
	return err
}

type GameSocket struct {
	rooms map[string]*GameRoom
}

func (g *GameSocket) UpdateGameState(ev game.Event, gm *game.Game) error {
	var errs = make([]error, 0)
	if room, ok := g.rooms[gm.ID.String()]; ok {
		for pi, player := range room.players {
			// send event and data to websocket
			sessions := func() []*game.Session {
				s := make([]*game.Session, len(gm.PlayerSessions))
				i := 0
				for _, sess := range gm.PlayerSessions {
					s[i] = sess
					i++
				}
				return s
			}()
			err := player.WriteJSON(map[string]any{
				"event": ev,
				"data": WSGameDto{
					ID:       &gm.ID,
					Settings: &gm.Settings,
					Sessions: sessions,
				},
			})
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					_ = player.Close()
					room.removePlayerAt(pi)
				}
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		strs := func(errs []error) (t []string) {
			for _, err := range errs {
				t = append(t, err.Error())
			}
			return t
		}(errs)
		return errors.New(strings.Join(strs, "\n"))
	}

	return nil
}

// tidyRooms removes rooms that have no players
// connected to them at intervals to prevent memory leaks
func (g *GameSocket) tidyRooms(duration time.Duration) {
	for {
		time.Sleep(duration)
		for id, room := range g.rooms {
			if len(room.players) == 0 {
				delete(g.rooms, id)
			}
		}
	}
}

// ServeHTTP handles websocket requests for GameSocket
// and connects the user to a game session if correct id is produced
func (g *GameSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("error upgrading to websockets", err)
		return
	}
	v := mux.Vars(r)
	if id, ok := v["id"]; !ok {
		fmt.Println("error getting game id")
		return
	} else {
		if room, ok := g.rooms[id]; ok {
			room.players = append(room.players, conn)
		} else {
			g.rooms[id] = &GameRoom{id, []*websocket.Conn{conn}}
		}
	}
}

func NewGameSocket() *GameSocket {
	sock := &GameSocket{make(map[string]*GameRoom)}
	go sock.tidyRooms(time.Hour)
	return sock
}

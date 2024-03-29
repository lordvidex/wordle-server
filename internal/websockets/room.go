package websockets

import (
	"github.com/lordvidex/wordle-wf/internal/game"
	"sync"
)

type Room struct {
	ID            string
	players       map[*Client]bool
	broadcast     chan interface{}
	join          chan *Client
	leave         chan *Client
	settings      game.Settings
	hasActiveGame bool
	owner         string
	mu            sync.Mutex
}

// NewRoom creates a new room for gamers playing game.Game with Room.ID
// and initializes all room's channels
func NewRoom(id string, settings game.Settings) *Room {
	return &Room{
		ID:            id,
		players:       make(map[*Client]bool),
		broadcast:     make(chan interface{}),
		join:          make(chan *Client),
		leave:         make(chan *Client),
		settings:      settings,
		hasActiveGame: false,
		owner:         "",
	}
}

func (r *Room) Close() error {
	var err error
	for player := range r.players {
		err = player.conn.Close()
		delete(r.players, player)
	}
	close(r.leave)
	close(r.join)
	close(r.broadcast)
	return err
}

func (r *Room) checkAssignOwner() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for player := range r.players {
		// owner exists
		if r.owner == player.playerID {
			return
		}
	}
	// no owner exists, assign one
	for player := range r.players {
		r.owner = player.playerID
		r.broadcast <- OwnerAssignedPayload(player.playerID, player.playerName)
		break
	}
}

func (r *Room) Run() {
	for {
		select {
		// a client joined the room
		case client := <-r.join:
			r.players[client] = true
			r.checkAssignOwner()
			r.broadcast <- &WSPayload{
				Event: EventPlayerJoined,
				Data: map[string]string{
					"player_id":   client.playerID,
					"player_name": client.playerName,
				},
			}

		// a client left the room
		case client := <-r.leave:
			delete(r.players, client)
			r.broadcast <- &WSPayload{
				// Event: game.EventPlayerLeft,
			}

		// broadcast message to all clients
		case msg := <-r.broadcast:
			for player := range r.players {
				select {
				case player.send <- msg:
				default:
					delete(r.players, player)
					close(player.send)
				}
			}
		}
	}
}

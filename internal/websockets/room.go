package websockets

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	ErrPayloadNotBytes = errors.New("payload is not bytes")
	ErrEmptyPayload    = errors.New("payload is empty")
	ErrGameInProgress  = errors.New("game in progress")
	ErrNoActiveGame    = errors.New("no active game")
)

type Room struct {
	ID                string
	players           map[*Client]bool
	broadcast         chan interface{}
	join              chan *Client
	leave             chan *Client
	settings          game.Settings
	activeGame        *game.Game
	owner             string
	mu                sync.Mutex
	createGameHandler game.CreateGameHandler
}

// NewRoom creates a new room for gamers playing game.Game with Room.ID
// and initializes all room's channels
func NewRoom(id string, settings game.Settings, cgh game.CreateGameHandler) *Room {
	return &Room{
		ID:                id,
		players:           make(map[*Client]bool),
		broadcast:         make(chan interface{}),
		join:              make(chan *Client),
		leave:             make(chan *Client),
		settings:          settings,
		createGameHandler: cgh,
	}
}

func (r *Room) HasActiveGame() bool {
	return r.activeGame != nil
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
			r.broadcast <- JoinPayload(client.playerID, client.playerName)
			// TODO: send room data to client

		// a client left the room
		case client := <-r.leave:
			delete(r.players, client)
			r.broadcast <- LeavePayload(client.playerID, client.playerName)

		// broadcast message to all clients
		case msg := <-r.broadcast:
			payload, err := unmarshalPayload(msg)
			if err != nil {
				logrus.WithError(err).Debug("error while converting message to payload")
				break
			}
			r.processPayload(&payload)
		}
	}
}

func (r *Room) processPayload(payload *WSPayload) {
	switch payload.Event {
	case EventGameStarted:
		err := r.StartGame()
		if err != nil {
			// send message to user about start failure
			for player := range r.players {
				if player.playerID == r.owner {
					player.send <- StartGameFailedPayload(err)
					break
				}
			}
			return
		}
		r.sendToAll(payload)
	case EventPlayerGuessed:
		if !r.HasActiveGame() {
			return // no active game
		}
		if r.settings.ViewOpponentsSessions {
			r.sendToAll(payload)
		}
	case EventPlayerJoined, EventPlayerLeft, EventOwnerAssigned:
		r.sendToAll(payload)
	}
}
func (r *Room) EndGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.HasActiveGame() {
		return ErrNoActiveGame
	}
	r.activeGame = nil
	return nil
}

func (r *Room) StartGame() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.HasActiveGame() {
		return ErrGameInProgress
	}
	// make call to create a new game
	ownerID, err := uuid.Parse(r.owner)
	if err != nil {
		return err
	}
	gm, err := r.createGameHandler.Handle(game.CreateGameCommand{
		CaptainID:   ownerID,
		Settings:    r.settings,
		InviteID:    r.ID,
		PlayerCount: len(r.players),
	})
	if err != nil {
		return err
	}
	r.activeGame = gm
	return nil
}

func (r *Room) sendToAll(payload *WSPayload) {
	for player := range r.players {
		select {
		case player.send <- payload:
		default:
			close(player.send)
		}
	}
}

func unmarshalPayload(msg interface{}) (WSPayload, error) {
	var payload WSPayload
	switch msg := msg.(type) {
	case []byte:
		err := json.Unmarshal(msg, &payload)
		if (WSPayload{}) == payload || err != nil {
			return payload, ErrEmptyPayload
		}
		return payload, err
	default:
		return payload, ErrPayloadNotBytes
	}
}

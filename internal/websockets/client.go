package websockets

import (
	"github.com/gorilla/websocket"
)

// Client is a websocket client and channel for game events
type Client struct {
	conn       *websocket.Conn
	send       chan interface{}
	room       *Room
	playerID   string
	playerName string
}

// NewClient creates a new websocket client
// and adds him to the room by sending him to the room's join channel
//
func NewClient(room *Room, conn *websocket.Conn, playerID string, playerName string) *Client {
	// create client
	cl := &Client{
		send:       make(chan interface{}),
		room:       room,
		conn:       conn,
		playerID:   playerID,
		playerName: playerName,
	}
	// start reading from the client and writing to the client
	go cl.ReadLoop()
	go cl.WriteLoop()

	// join room
	room.join <- cl
	return cl
}

func (c *Client) ReadLoop() {
	defer func() {
		_ = c.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		c.room.broadcast <- msg
	}
}

func (c *Client) Close() error {
	c.room.leave <- c
	close(c.send)
	return c.conn.Close()
}

func (c *Client) WriteLoop() {
	defer func() {
		_ = c.Close()
	}()
	for msg := range c.send {
		err := c.conn.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}

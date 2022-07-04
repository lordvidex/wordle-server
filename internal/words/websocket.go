package words

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
)

type WebsocketHandler struct {
}

func (w WebsocketHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		conn.Close()
	}
	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
		}

		fmt.Println(b)
		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Server is returning %s", b)))
		if err != nil {
			fmt.Println("Error writing message")
		}
		err = conn.WriteJSON(map[string]string{
			"title":   "Hello World",
			"message": string(b),
		})
		if err != nil {
			fmt.Println("An error occurred with ", err)
		}
	}
}

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{}
}

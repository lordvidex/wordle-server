package words

//import "golang.org/x/net/websocket"
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
		fmt.Printf("error occured while upgrading the connection: %s", err)
	}

	go func() {
		for {
			_, b, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err) {
					break
				}
				fmt.Println("An error occured reading message", err)
			}

			fmt.Println(b)
			err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Server is returning %s", b)))
			if err != nil {
				fmt.Println("Error writing message")
			}

		}
	}()
	err = conn.WriteJSON(map[string]string{"message": "hello"})

	if err != nil {
		fmt.Println("Error writing to connection")
	}
}

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{}
}

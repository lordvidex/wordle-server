package game

type Lobby interface {
	CreateRoom(ID string) error
	JoinRoom(ID string) error
	PlayInRoom(ID string, word string) error
	FindRoom(ID string) (Room, error)
}

package websockets

// Event is an enum for all the events that Listeners need to subscribe to
// when listening to updates to the Game state
type Event string

const (
	EventPlayerJoined  Event = "PlayerJoined"
	EventPlayerLeft    Event = "PlayerLeft"
	EventGameStarted   Event = "GameStarted"
	EventGameEnded     Event = "GameEnded"
	EventPlayerGuessed Event = "PlayerGuessed"
)

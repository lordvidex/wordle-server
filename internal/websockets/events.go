package websockets

import "github.com/lordvidex/wordle-wf/internal/game"

// Event is an enum for all the events that Listeners need to subscribe to
// when listening to updates to the Game state
type Event string

const (
	EventPlayerJoined    Event = "PlayerJoined"
	EventPlayerLeft      Event = "PlayerLeft"
	EventOwnerAssigned   Event = "OwnerAssigned"
	EventGameStarted     Event = "GameStarted"
	EventGameStartFailed Event = "GameStartFailed"
	EventGameEnded       Event = "GameEnded"
	EventPlayerGuessed   Event = "PlayerGuessed"
	//EventRoomData        Event = "RoomData"
)

func OwnerAssignedPayload(playerID string, playerName string) *WSPayload {
	return &WSPayload{
		Event: EventOwnerAssigned,
		Data: map[string]string{
			"player_id":   playerID,
			"player_name": playerName,
		},
	}
}

func StartGameFailedPayload(err error) *WSPayload {
	return &WSPayload{
		Event: EventGameStartFailed,
		Data: map[string]string{
			"error": err.Error(),
		},
	}
}

func JoinPayload(playerID string, playerName string) *WSPayload {
	return &WSPayload{
		Event: EventPlayerJoined,
		Data: map[string]string{
			"player_id":   playerID,
			"player_name": playerName,
		},
	}
}

func LeavePayload(playerID string, playerName string) *WSPayload {
	return &WSPayload{
		Event: EventPlayerLeft,
		Data: map[string]string{
			"player_id":   playerID,
			"player_name": playerName,
		},
	}
}

func StartPayload(gameData game.Game) *WSPayload {
	return &WSPayload{
		Event: EventGameStarted,
		Data:  gameData,
	}
}

func EndPayload(gameData game.Game) *WSPayload {
	return &WSPayload{
		Event: EventGameEnded,
		Data:  gameData,
	}
}

func GuessPayload(playerID string, playerName string) *WSPayload {
	return &WSPayload{
		Event: EventPlayerGuessed,
		Data: map[string]string{
			"player_id":   playerID,
			"player_name": playerName,
		},
	}
}

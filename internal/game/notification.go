package game

// NotificationService is responsible for dispatching updates to the game state.
type NotificationService interface {
	UpdateGameState(ev Event, game *Game) error
}

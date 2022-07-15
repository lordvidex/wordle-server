package game

// Repository is the interface for storing game properties and state
type Repository interface {
	// Create creates a new game and the settings at the same time
	// by default, the game appends to all related tables
	// returns the game instance back newly from the database
	Create(game *Game) (*Game, error)

	// Join adds a player to a game if he provides the correct id
	Join(gameId string, player *Player) error

	// UpdateSettings updates the settings of a game that is yet to be started
	UpdateSettings(settings *Settings, gameID string) error

	// Start starts a game and begins monitoring user's sessions
	Start(gameID string) error

	// Find returns a game by its ID only from it's table UNLESS
	// eager is provided.
	// eager should be interfaces of the models that are to be eagerly loaded
	FindByID(gameId string, eager ...interface{}) (*Game, error)

	// FindByInviteID returns a game if the invite id provided matches any of the
	// available games in the database that has started but has not ended
	FindByInviteID(inviteId string, eager ...interface{}) ([]Game, error)

	// FindAll returns all games from the database
	FindAll(eager ...interface{}) ([]Game, error)

	// Delete deletes a game from the database with CASCADE on Settings, Word and Sessions
	Delete(gameId string) error
}

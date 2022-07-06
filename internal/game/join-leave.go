package game

type JoinGameCommand struct {
	Player *Player
	GameID string
}

// RandomJoinGameCommand is a command to join a game randomly
// without a Game ID. A Player proposes a game Settings and checks if any user with that
// setting exists in the buffer. Else, he joins the buffer and wait for another player to join
// until he times out.
type RandomJoinGameCommand struct {
	Player   *Player
	Settings *Settings
}

type LeaveGameCommand struct {
	Player *Player
	GameID string
}

type JoinGameCommandHandler interface {
	HandleJoin(command JoinGameCommand) error
	HandleRandomJoin(command RandomJoinGameCommand) error
}

type LeaveGameCommandHandler interface {
	HandleLeave(command LeaveGameCommand) error
}

type joinGameCommandHandler struct {
	repo     Repository
	notifier NotificationService
}

func NewJoinGameCommandHandler(repo Repository, notifier NotificationService) JoinGameCommandHandler {
	return &joinGameCommandHandler{repo, notifier}
}

func (c *joinGameCommandHandler) HandleJoin(command JoinGameCommand) error {
	game, err := c.repo.Find(command.GameID)
	if err != nil {
		return err
	}
	// create a new session for player
	game.PlayerSessions[*command.Player] = NewSession(command.Player)
	game, err = c.repo.Save(game)
	if err != nil {
		return err
	}
	// notify all players of the new player
	return c.notifier.UpdateGameState(EventPlayerJoined, game)
}

func (c *joinGameCommandHandler) HandleRandomJoin(command RandomJoinGameCommand) error {

}

func (c *joinGameCommandHandler) HandleLeave(command LeaveGameCommand) error {
	game, err := c.repo.Find(command.GameID)
	if err != nil {
		return err
	}
	delete(game.PlayerSessions, *command.Player)
	game, err = c.repo.Save(game)
	if err != nil {
		return err
	}
	return c.notifier.UpdateGameState(EventPlayerLeft, game)
}

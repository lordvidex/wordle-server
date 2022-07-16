package pg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lordvidex/wordle-wf/internal/auth"
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
	"github.com/lordvidex/wordle-wf/internal/db/pg/mapper"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/lordvidex/wordle-wf/internal/words"
	"time"
)

type gameRepository struct {
	*pg.Queries
	c     context.Context
	pgxDB *pgxpool.Pool
}

// external errors
var (
	ErrUUIDParse           = errors.New("uuid parse error")
	ErrUserNotFound        = errors.New("user not found")
	ErrDataNotFound        = errors.New("data not found")
	ErrJoinGame            = errors.New("join game error")
	ErrStartGame           = errors.New("start game error")
	ErrFetchingPlayerGuess = errors.New("fetching words for player error")
)

func NewGameRepository(db *pgxpool.Pool) game.Repository {
	return &gameRepository{
		c:       context.Background(),
		Queries: pg.New(db),
		pgxDB:   db,
	}
}

func (g *gameRepository) Create(game *game.Game) (*game.Game, error) {
	tx, err := g.pgxDB.BeginTx(g.c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(g.c)
	}()
	// create game
	err = g.WithTx(tx).CreateGame(g.c, pg.CreateGameParams{ID: game.ID, InviteID: game.InviteID})
	if err != nil {
		return nil, err
	}
	// create settings
	err = g.WithTx(tx).CreateGameSettings(g.c, pg.CreateGameSettingsParams{
		GameID:                   uuid.NullUUID{UUID: game.ID, Valid: true},
		WordLength:               sql.NullInt16{Int16: int16(game.Settings.WordLength), Valid: true},
		Trials:                   sql.NullInt16{Int16: int16(game.Settings.Trials), Valid: true},
		PlayerCount:              sql.NullInt16{Int16: int16(game.Settings.PlayerCount), Valid: true},
		HasAnalytics:             sql.NullBool{Bool: game.Settings.Analytics, Valid: true},
		ShouldRecordTime:         sql.NullBool{Bool: game.Settings.RecordTime, Valid: true},
		CanViewOpponentsSessions: sql.NullBool{Bool: game.Settings.ViewOpponentsSessions, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	err = tx.Commit(g.c)
	if err != nil {
		return nil, err
	}
	return game, nil
}

func (g *gameRepository) Join(gameId string, player *game.Player) error {
	gameUUID, err := uuid.Parse(gameId)
	if err != nil {
		return ErrUUIDParse
	}
	if player.User == nil {
		return ErrUserNotFound
	}
	if player.Name == "" {
		player.Name = player.User.Name
	}
	_, err = g.Queries.CreateGamePlayer(g.c, pg.CreateGamePlayerParams{
		GameID: gameUUID,
		Name:   player.Name,
		UserID: player.User.ID,
	})
	if err != nil {
		return ErrJoinGame
	}
	return nil
}

func (g *gameRepository) UpdateSettings(settings *game.Settings, gameID string) error {
	gameUUID, err := uuid.Parse(gameID)
	if err != nil {
		return ErrUUIDParse
	}
	if settings == nil {
		return ErrDataNotFound
	}
	_, err = g.UpdateGameSettings(g.c, pg.UpdateGameSettingsParams{
		GameID:                   uuid.NullUUID{UUID: gameUUID, Valid: true},
		WordLength:               sql.NullInt16{Int16: int16(settings.WordLength), Valid: settings.WordLength != 0},
		Trials:                   sql.NullInt16{Int16: int16(settings.Trials), Valid: settings.Trials != 0},
		PlayerCount:              sql.NullInt16{Int16: int16(settings.PlayerCount), Valid: settings.PlayerCount != 0},
		HasAnalytics:             sql.NullBool{Bool: settings.Analytics, Valid: true},
		ShouldRecordTime:         sql.NullBool{Bool: settings.RecordTime, Valid: true},
		CanViewOpponentsSessions: sql.NullBool{Bool: settings.ViewOpponentsSessions, Valid: true},
	})
	return err
}

func (g *gameRepository) Start(gameID string) error {
	gameUUID, err := uuid.Parse(gameID)
	if err != nil {
		return ErrUUIDParse
	}
	err = g.StartGame(g.c, pg.StartGameParams{
		StartTime: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        gameUUID,
	})
	if err != nil {
		return ErrStartGame
	}
	return nil
}

func (g *gameRepository) FindByID(gameId string, eager ...interface{}) (*game.Game, error) {
	var gm game.Game
	// parse uuid
	gameUUID, err := uuid.Parse(gameId)
	if err != nil {
		return nil, ErrUUIDParse
	}
	//
	//// fetch player if eager Player
	//if isEager([2]interface{}{game.Player{}, &game.Player{}}, eager) {
	//	playersResult = make(chan ResultsAndError[*game.Player])
	//	go func() {
	//		defer close(playersResult)
	//		result, err := g.Queries.GetPlayersInGame(g.c, gameUUID)
	//		if err != nil {
	//			playersResult <- ResultsAndError[*game.Player]{Err: err}
	//		}
	//		playersResult <- ResultsAndError[*game.Player]{Data: mapper.GetPlayersInGame(result)}
	//	}()
	//}
	//
	// fetch game
	gameRow, err := g.FindById(g.c, gameUUID)
	if err != nil {
		return nil, err
	}
	mapper.FindByIdRow(gameRow, &gm)
	//// receive players data
	//if playersResult != nil {
	//	chanResponse := <-playersResult
	//	// handle errors
	//	if chanResponse.Err != nil {
	//		return nil, chanResponse.Err
	//	}
	//	players = chanResponse.Data
	//}
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//// assemble data

	//
	//// add players
	//for _, p := range players.Data {
	//	gm.Sessions[*p] = nil
	//}
	return &gm, nil
}

func (g *gameRepository) FindByInviteID(inviteId string, eager ...interface{}) ([]game.Game, error) {
	rows, err := g.FindByInviteId(g.c, sql.NullString{String: inviteId, Valid: true})
	if err != nil {
		return nil, err
	}
	// get players in goroutine
	pipe := make(chan game.Game)
	defer close(pipe)
	// get words played by player
	sessionPipe := make(chan *game.Session)
	defer close(sessionPipe)

	gms := make([]game.Game, len(rows))

	// loop through each game
	for _, row := range rows {
		gm := game.Game{
			ID:       row.ID,
			InviteID: row.InviteID,
			Settings: game.Settings{
				WordLength:            int(row.WordLength.Int16),
				Trials:                int(row.Trials.Int16),
				PlayerCount:           int(row.PlayerCount.Int16),
				Analytics:             row.HasAnalytics.Bool,
				RecordTime:            row.ShouldRecordTime.Bool,
				ViewOpponentsSessions: row.CanViewOpponentsSessions.Bool,
			},
		}
		// worker to get players in each game
		fetchPlayerWords := false
		if isEager([2]interface{}{words.Word{}, &words.Word{}}, eager) {
			fetchPlayerWords = true
		}
		go g.playersInGameWorker(&gm, pipe, fetchPlayerWords, sessionPipe)
	}
	// assemble
	for i := range rows {
		gms[i] = <-pipe
	}
	return gms, nil
}

func (g *gameRepository) FindAll(eager ...interface{}) ([]game.Game, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) Delete(gameId string) error {
	gameUUID, err := uuid.Parse(gameId)
	if err != nil {
		return ErrUUIDParse
	}
	return g.Queries.DeleteGame(g.c, gameUUID)
}

// worker functions
//

// playersInGameWorker fetches players in game and sends the result to the channel
func (g *gameRepository) playersInGameWorker(gm *game.Game,
	pipe chan<- game.Game,
	fetchWords bool,
	sessionPipe chan<- *game.Session,
) {
	list, _ := g.GetPlayersInGame(g.c, gm.ID)
	gm.Sessions = make([]*game.Session, len(list))
	for i, p := range list {
		gm.Sessions[i] = &game.Session{
			Player: &game.Player{
				ID:   p.ID,
				Name: p.Name,
				User: &auth.User{
					ID:    p.UserID,
					Name:  p.UserName.String,
					Email: p.Email.String,
				}},
		}
		if fetchWords {
			go g.wordsByPlayerWorker(gm.Sessions[i], sessionPipe)
		}
	}
	pipe <- *gm
}

// wordsByPlayerWorker fetches words for a player and sends the result to the channel
// WARNING: this function discards errors when getting player words
func (g *gameRepository) wordsByPlayerWorker(sess *game.Session, pipe chan<- *game.Session) {
	w, err := g.WordsPlayedBy(g.c, sess.Player.ID)
	if err != nil {
		goto end
	}
	sess.Guesses, err = mapper.WordsPlayedBy(w)
	if err != nil {
		goto end
	}
end:
	pipe <- sess
}

// helper funcs

// isEager returns true if the interface t should be eagerly loaded (contained in ints)
// for this function to properly work, `t` and all interface in `ints` must be zero struct types
// e.g. `X{}`
// i.e. without values
// t takes an array so that pointer as well as struct values can be checked
func isEager(t [2]interface{}, ints ...interface{}) bool {
	for _, i := range ints {
		if i != nil && (i == t[0] || i == t[1]) {
			return true
		}
	}
	return false
}

type ResultsAndError[T any] struct {
	Data []T
	Err  error
}

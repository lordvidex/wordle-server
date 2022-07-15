package pg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
	"github.com/lordvidex/wordle-wf/internal/db/pg/mapper"
	"github.com/lordvidex/wordle-wf/internal/game"
	"time"
)

type gameRepository struct {
	*pg.Queries
	c     context.Context
	pgxDB *pgx.Conn
}

// external errors
var (
	ErrUUIDParse    = errors.New("uuid parse error")
	ErrUserNotFound = errors.New("user not found")
	ErrDataNotFound = errors.New("data not found")
	ErrJoinGame     = errors.New("join game error")
	ErrStartGame    = errors.New("start game error")
)

// internal errors
var ()

func NewGameRepository(db *pgx.Conn) game.Repository {

	return &gameRepository{
		c:       context.Background(),
		Queries: pg.New(db),
		pgxDB:   db,
	}
}

func (g *gameRepository) Create(game *game.Game) (*game.Game, error) {
	tx, err := g.pgxDB.Begin(g.c)
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
	// create transaction
	tx, err := g.pgxDB.Begin(g.c)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(g.c)
	}()
	var gm game.Game
	//var players = make([]*game.Session, 0)

	// parse uuid
	gameUUID, err := uuid.Parse(gameId)
	if err != nil {
		return nil, ErrUUIDParse
	}
	//
	//// fetch player if eager Player
	//if IsEager([2]interface{}{game.Player{}, &game.Player{}}, eager) {
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
	gameRow, err := g.WithTx(tx).FindById(g.c, gameUUID)

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

	mapper.FindByIdRow(gameRow, &gm)
	//
	//// add players
	//for _, p := range players.Data {
	//	gm.PlayerSessions[*p] = nil
	//}
	err = tx.Commit(g.c)
	if err != nil {
		return nil, err
	}
	return &gm, nil
}

func (g *gameRepository) FindByInviteID(inviteId string, eager ...interface{}) ([]game.Game, error) {
	//TODO implement me
	panic("implement me")
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

// helper funcs

// IsEager returns true if the interface t should be eagerly loaded (contained in ints)
// for this function to properly work, `t` and all interface in `ints` must be zero struct types
// e.g. `X{}`
// i.e. without values
// t takes an array so that pointer as well as struct values can be checked
func IsEager(t [2]interface{}, ints ...interface{}) bool {
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

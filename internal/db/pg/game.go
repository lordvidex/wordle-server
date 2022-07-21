package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lordvidex/wordle-wf/internal/words"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
	"github.com/lordvidex/wordle-wf/internal/db/pg/mapper"
	"github.com/lordvidex/wordle-wf/internal/game"
)

type gameRepository struct {
	*pg.Queries
	c     context.Context
	pgxDB *pgx.Conn
}

// external errors
var (
	ErrUUIDParse    = errors.New("uuid parse error")
	ErrDataNotFound = errors.New("data not found")
)

func NewGameRepository(db *pgx.Conn) game.Repository {
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
	err = g.WithTx(tx).CreateGame(g.c, pg.CreateGameParams{ID: game.ID,
		InviteID:    game.InviteID,
		Word:        game.Word.String(),
		PlayerCount: int16(game.PlayerCount),
		StartTime:   game.StartTime,
	})
	if err != nil {
		return nil, err
	}
	// create settings
	err = g.WithTx(tx).CreateGameSettings(g.c, pg.CreateGameSettingsParams{
		GameID:                   uuid.NullUUID{UUID: game.ID, Valid: true},
		WordLength:               sql.NullInt16{Int16: int16(game.Settings.WordLength), Valid: true},
		Trials:                   sql.NullInt16{Int16: int16(game.Settings.Trials), Valid: true},
		MaxPlayerCount:           sql.NullInt16{Int16: int16(game.Settings.MaxPlayerCount), Valid: true},
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

func (g *gameRepository) Delete(id uuid.UUID) error {
	return g.Queries.DeleteGame(g.c, id)
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
		MaxPlayerCount:           sql.NullInt16{Int16: int16(settings.MaxPlayerCount), Valid: settings.MaxPlayerCount != 0},
		HasAnalytics:             sql.NullBool{Bool: settings.Analytics, Valid: true},
		ShouldRecordTime:         sql.NullBool{Bool: settings.RecordTime, Valid: true},
		CanViewOpponentsSessions: sql.NullBool{Bool: settings.ViewOpponentsSessions, Valid: true},
	})
	return err
}

func (g *gameRepository) FindByID(gameID uuid.UUID) (*game.Game, error) {
	gameRow, err := g.FindById(g.c, gameID)
	if err != nil {
		return nil, err
	}
	gm := mapper.FindByIDRow(gameRow)

	// fetch players in game
	players, err := g.Queries.GetPlayersResultInGame(g.c, uuid.NullUUID{UUID: gm.ID})
	if err != nil {
		return nil, err
	}
	gm.Sessions = mapper.GetPlayersResultInGame(players)
	for _, s := range gm.Sessions {
		// fetch words in game
		playerWords, err := g.Queries.PlayerWordsInGame(g.c, pg.PlayerWordsInGameParams{
			PlayerID: uuid.NullUUID{
				UUID:  s.Player.ID,
				Valid: true,
			},
			GameID: uuid.NullUUID{
				UUID:  gm.ID,
				Valid: true,
			},
		})
		if err != nil {
			return nil, err
		}
		s.Guesses = mapper.PlayerWordsInGame(playerWords)
	}
	// map players to game.Players
	return gm, nil
}

func (g *gameRepository) UpdateGameResult(
	gameID uuid.UUID,
	playerID uuid.UUID,
	playerName string,
	points func(position int) int,
	words []*words.Word,
) error {
	tx, err := g.pgxDB.BeginTx(g.c, pgx.TxOptions{})
	if err != nil {
		return err
	}
	lastPosition, err := g.Queries.
		WithTx(tx).
		GetPlayersResultCountInGame(g.c, uuid.NullUUID{UUID: gameID})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(g.c)
	}()

	myPosition := int32(lastPosition + 1)
	myPoint := points(int(myPosition))
	id, err := g.Queries.WithTx(tx).UpdateGameResult(g.c, pg.UpdateGameResultParams{
		GameID:       uuid.NullUUID{UUID: gameID, Valid: true},
		PlayerID:     uuid.NullUUID{UUID: playerID, Valid: true},
		UserGameName: playerName,
		Points: sql.NullInt32{
			Int32: int32(myPoint),
			Valid: true,
		},
		Position: sql.NullInt32{
			Int32: myPosition,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}
	// update the words played
	var sqlWords = mapper.InsertWords(uuid.NullUUID{UUID: id, Valid: true}, words)
	_, err = g.Queries.WithTx(tx).InsertWords(g.c, sqlWords)
	if err != nil {
		return err
	}
	return tx.Commit(g.c)
}

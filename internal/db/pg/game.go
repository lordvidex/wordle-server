package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/wordle-wf/internal/db/pg/gen"
	"github.com/lordvidex/wordle-wf/internal/game"
)

type gameRepository struct {
	c     context.Context
	db    *pg.Queries
	pgxDB *pgx.Conn
}

func NewGameRepository(db *pgx.Conn) game.Repository {
	return &gameRepository{
		c:     context.Background(),
		db:    pg.New(db),
		pgxDB: db,
	}
}

func (g *gameRepository) Create(game *game.Game) (*game.Game, error) {
	var lettersJSON pgtype.JSON
	err := (&lettersJSON).Set(game.Word.Letters)
	if err != nil {
		return nil, err
	}
	// start a transaction
	tx, err := g.pgxDB.Begin(g.c)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(g.c)

	// create the word for the game
	word, err := g.db.WithTx(tx).CreateWord(g.c, pg.CreateWordParams{
		Wordid:      game.Word.ID,
		Timeplayed:  game.Word.TimePlayed,
		Lettersjson: lettersJSON,
	})
	if err != nil {
		return nil, err
	}

	// create the game
	err = g.db.WithTx(tx).CreateGame(g.c, pg.CreateGameParams{
		ID:     game.ID,
		WordID: uuid.NullUUID{UUID: word.ID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// create the game settings
	err = g.db.WithTx(tx).CreateGameSettings(g.c, pg.CreateGameSettingsParams{
		GameID: uuid.NullUUID{UUID: game.ID, Valid: true},
		WordLength: sql.NullInt16{
			Int16: int16(game.Settings.WordLength),
			Valid: game.Settings.WordLength > 0,
		},
		Trials: sql.NullInt16{
			Int16: int16(game.Settings.Trials),
			Valid: game.Settings.Trials != 0,
		},
		PlayerCount: sql.NullInt16{
			Int16: int16(game.Settings.PlayerCount),
			Valid: game.Settings.PlayerCount != 0,
		},
		HasAnalytics: sql.NullBool{
			Bool:  game.Settings.Analytics,
			Valid: true,
		},
		ShouldRecordTime: sql.NullBool{
			Bool:  game.Settings.RecordTime,
			Valid: true,
		},
		CanViewOpponentsSessions: sql.NullBool{
			Bool:  game.Settings.ViewOpponentsSessions,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// create the player sessions for the game
	for player, session := range game.PlayerSessions {
		sessionID, _ := g.db.WithTx(tx).AddPlayerSessionToGame(g.c, pg.AddPlayerSessionToGameParams{
			GameID:   game.ID,
			PlayerID: player.ID,
		})
		// add the word guesses for each player session
		guesses := make([]pg.AddPlayerGuessParams, len(session.Guesses))
		for i, guess := range session.Guesses {
			guesses[i] = pg.AddPlayerGuessParams{
				GameSessionID: uuid.NullUUID{
					UUID:  sessionID,
					Valid: true,
				},
				WordID: uuid.NullUUID{
					UUID:  guess.ID,
					Valid: true,
				},
			}
		}
		cnt, err := g.db.WithTx(tx).AddPlayerGuess(g.c, guesses)
		if err != nil {
			return nil, err
		}
		if cnt <= 0 {
			fmt.Println("warning.. no player guesses added to session")
		}
	}
	return game, nil
}

func (g *gameRepository) Save(game *game.Game) (*game.Game, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) Find(id string) (*game.Game, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) FindAll() ([]game.Game, error) {
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) Close() error {
	//TODO implement me
	panic("implement me")
}

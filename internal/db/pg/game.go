package pg

import (
	"context"
	"github.com/jackc/pgx/v4"
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
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
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) Join(gameId string, player *game.Player) error {
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) UpdateSettings(settings *game.Settings, gameID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) Start(gameID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *gameRepository) FindByID(gameId string, eager ...interface{}) (*game.Game, error) {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

package pg

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type GameModel struct {
	ID   uuid.UUID `bun:",pk,type:uuid"`
	Word WordModel `bun:"rel:belongs-to,join:word_id=id"`
}
type WordModel struct {
	ID uuid.UUID `bun:",pk,type:uuid"`
}

func (*GameModel) FromGame(g *game.Game) *GameModel {
	return &GameModel{
		ID: g.ID,
		Word: WordModel{
			ID: g.Word.ID,
		},
	}
}

type bunRepository struct {
	db *bun.DB
	c  context.Context
}

func NewBunRepository(db *sql.DB) game.Repository {
	bunDB := bun.NewDB(db, pgdialect.New())
	return &bunRepository{db: bunDB, c: context.Background()}
}

func (g *bunRepository) Create(gm *game.Game) (*game.Game, error) {
	var newGm game.Game
	_, err := g.db.NewInsert().Model(gm).Exec(g.c, &newGm)
	return &newGm, err
}

func (g *bunRepository) Close() error {
	return g.db.Close()
}

func (g *bunRepository) Save(game *game.Game) (*game.Game, error) {
	_, err := g.db.NewUpdate().Model(game).Returning("").Exec(context.Background())
	return game, err
}

func (g *bunRepository) Find(id string) (*game.Game, error) {
	var gm game.Game
	err := g.db.NewSelect().Where("id = ?", uuid.MustParse(id)).Limit(1).Scan(context.Background(), &gm)
	return &gm, err
}

func (g *bunRepository) FindAll() ([]game.Game, error) {
	//TODO implement me
	panic("implement me")
}

func (g *bunRepository) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

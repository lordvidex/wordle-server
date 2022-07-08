package mongo

import (
	"context"
	"github.com/lordvidex/wordle-wf/internal/game"
	"go.mongodb.org/mongo-driver/mongo"
)

// GameModel represents a game.Game model in mongo database
type GameModel struct {
	ID string `bson:"_id"`
}

type GameRepository struct {
	db *mongo.Database
}

func (g *GameRepository) Create(game *game.Game) (*game.Game, error) {
	_, err := g.db.Collection("games").InsertOne(context.TODO(), game)
	return game, err
}

func (g *GameRepository) Save(game *game.Game) (*game.Game, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GameRepository) Find(id string) (*game.Game, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GameRepository) FindAll() ([]game.Game, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GameRepository) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GameRepository) Close() error {
	//TODO implement me
	panic("implement me")
}

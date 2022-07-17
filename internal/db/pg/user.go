package pg

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/auth"
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
	"github.com/lordvidex/wordle-wf/internal/game"
)

type userRepository struct {
	Queries *pg.Queries
	pgxDB   *pgx.Conn
	c       context.Context
}

func NewUserRepository(db *pgx.Conn) auth.Repository {
	return &userRepository{pgxDB: db, Queries: pg.New(db), c: context.Background()}
}

func (u *userRepository) FindByEmail(email string) (*game.Player, error) {
	player, err := u.Queries.GetPlayerByEmail(u.c, email)
	if err != nil {
		return nil, err
	}
	return &game.Player{
		ID:       player.ID,
		Name:     player.Name,
		Email:    player.Email,
		Password: player.Password,
	}, nil
}

func (u *userRepository) FindByID(id uuid.UUID) (*game.Player, error) {
	user, err := u.Queries.GetPlayerByID(u.c, id)
	if err != nil {
		return nil, err
	}
	return &game.Player{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (u *userRepository) Create(name string, email string, password string) (*game.Player, error) {
	newUser, err := u.Queries.CreatePlayer(u.c, pg.CreatePlayerParams{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &game.Player{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}, nil
}

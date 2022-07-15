package pg

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/wordle-wf/internal/auth"
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
)

type userRepository struct {
	*pg.Queries
	pgxDB *pgx.Conn
	c     context.Context
}

func NewUserRepository(db *pgx.Conn) auth.Repository {
	return &userRepository{pgxDB: db, Queries: pg.New(db), c: context.Background()}
}

func (u *userRepository) FindByEmail(email string) (*auth.User, error) {
	user, err := u.Queries.GetUserByEmail(u.c, email)
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (u *userRepository) FindByID(id uuid.UUID) (*auth.User, error) {
	user, err := u.Queries.GetUserByID(u.c, id)
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (u *userRepository) Create(user *auth.User) (*auth.User, error) {
	//TODO implement me
	panic("implement me")
}

package pg

import (
	"context"
	"github.com/jackc/pgx/v4"

	"github.com/google/uuid"
	"github.com/lordvidex/wordle-wf/internal/auth"
	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
)

type userRepository struct {
	Queries *pg.Queries
	pgxDB   *pgx.Conn
	c       context.Context
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
	newUser, err := u.Queries.InsertUser(u.c, pg.InsertUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}
	return &auth.User{
		ID:    newUser.ID,
		Name:  newUser.Name,
		Email: newUser.Email,
	}, nil
}

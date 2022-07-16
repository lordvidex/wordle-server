package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"testing"

	pg "github.com/lordvidex/wordle-wf/internal/db/pg/gen"
)

const (
	testConnString = "postgres://postgres:postgres@localhost:5432/wordle_wf?sslmode=disable"
)

var mockUserRepo *userRepository
var mockGameRepo *gameRepository

func TestMain(m *testing.M) {
	conn, err := pgx.Connect(context.Background(), testConnString)
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	fmt.Println("connected to database")
	mockUserRepo = &userRepository{
		Queries: pg.New(conn),
		c:       context.Background(),
		pgxDB:   conn,
	}
	mockGameRepo = &gameRepository{
		Queries: pg.New(conn),
		c:       context.Background(),
		pgxDB:   conn,
	}
	m.Run()
}

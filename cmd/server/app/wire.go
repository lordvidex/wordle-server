//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/wordle-wf/internal/adapters"
	"github.com/lordvidex/wordle-wf/internal/auth"
	"github.com/lordvidex/wordle-wf/internal/db/pg"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/lordvidex/wordle-wf/internal/websockets"
	"github.com/lordvidex/wordle-wf/internal/words"
	"time"
)

func injectAuth(db *pgx.Conn, tokenSecret string, tokenIAT time.Duration) auth.UseCases {
	wire.Build(
		auth.NewUseCases,
		pg.NewUserRepository,
		adapters.NewPASETOTokenHelper,
		adapters.NewBcryptHelper,
	)
	return auth.UseCases{}
}

func injectGame(db *pgx.Conn, wrh words.RandomHandler, socket *websockets.GameSocket) game.UseCases {
	wire.Build(
		game.NewUseCases,
		pg.NewGameRepository,
		adapters.NewUniUriGenerator,
		adapters.NewAwardSystem,
		wire.Bind(new(game.LobbyCreator), new(*websockets.GameSocket)),
	)
	return game.UseCases{}
}

func injectRandomWordHandler(t words.UseCases) words.RandomHandler {
	return t.RandomWordHandler
}

func injectWord() words.UseCases {
	wire.Build(
		words.NewUseCases,
		adapters.NewLocalStringGenerator,
	)
	return words.UseCases{}
}

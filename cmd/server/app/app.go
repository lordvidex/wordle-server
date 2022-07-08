// Package app represents the entry point of the application.
// This is where dependencies are wired together and the application is started.
package app

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lordvidex/wordle-wf/internal/adapters"
	"github.com/lordvidex/wordle-wf/internal/auth"
	"github.com/lordvidex/wordle-wf/internal/db/pg"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/lordvidex/wordle-wf/internal/middleware"
	"github.com/lordvidex/wordle-wf/internal/websockets"
	"github.com/lordvidex/wordle-wf/internal/words"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"net/http"
)

func Start() {
	dsn := "postgres://postgres:@localhost:5432/test?sslmode=disable"
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	pgDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	router := mux.NewRouter()
	gameSocket := websockets.NewGameSocket()
	gameUsecase := game.NewUseCases(
		pg.NewBunRepository(pgDB),
		words.NewRandomHandler(adapters.NewLocalStringGenerator()),
		gameSocket)

	registerApi(router, gameUsecase)
	registerWS(router, gameSocket)
	registerAsset(router)
	printEndpoints(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func registerAsset(router *mux.Router) {
	router.Handle("/", http.FileServer(http.Dir("./resources")))
}

func registerWS(router *mux.Router, ws http.Handler) {
	router.Handle("/live", ws)
}

// registerApi registers the API endpoints.
func registerApi(router *mux.Router) {
	// main api endpoint
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.HandleError, middleware.JSONContent, middleware.Logger)

	// words endpoint
	wordsRouter := apiRouter.PathPrefix("/words").Subrouter()
	RegisterWordsHandler(wordsRouter)

	// auth endpoints
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	auth.RegisterHandler(authRouter)
}

// printEndpoints prints the endpoints that are exposed for api consumption
func printEndpoints(r *mux.Router) {
	if err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		fmt.Printf("%v %s\n", methods, path)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

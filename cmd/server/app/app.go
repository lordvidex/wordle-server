// Package app represents the entry point of the application.
// This is where dependencies are wired together and the application is started.
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/wordle-wf/internal/adapters"
	"github.com/lordvidex/wordle-wf/internal/auth"
	"github.com/lordvidex/wordle-wf/internal/db/pg"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/lordvidex/wordle-wf/internal/middleware"
	"github.com/lordvidex/wordle-wf/internal/websockets"
	"github.com/lordvidex/wordle-wf/internal/words"
	"github.com/spf13/viper"
)

type RouteBuilder struct {
	router *mux.Router
}

func (routeBuilder *RouteBuilder) MakeRoute(path string, f func(RouteBuilder, *mux.Router)) *RouteBuilder {
	apiRouter := routeBuilder.router.PathPrefix(path).Subrouter()
	f(RouteBuilder{router: apiRouter}, apiRouter)
	return routeBuilder
}

func connectDB(c *DBConfig) (*pgx.Conn, error) {
	var dsn string
	if c.Url != "" {
		dsn = c.Url
	} else {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, 5432, c.DBName)
	}
	pgConn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v\n", err)
	}
	err = pgConn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	log.Println("Connected to database")
	return pgConn, err
}

func Start() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatal("error occured loading configs", err)
	}
	// connect to database
	pgConn, err := connectDB(conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = pgConn.Close(context.Background())
	}()

	// repositories
	gameRepo := pg.NewGameRepository(pgConn)

	// services and dependencies
	var gameSocket *websockets.GameSocket
	defer func() {
		if err := gameSocket.Close(); err != nil {
			log.Println("error closing websocket", err)
		}
	}()

	// usecases and application layer components
	wordsUsecase := words.NewUseCases(
		adapters.NewLocalStringGenerator(),
		nil,
	)
	gameUsecase := game.NewUseCases(
		gameRepo,
		wordsUsecase.Queries.GetRandomWordHandler,
		gameSocket)

	// adapters and external services
	gameSocket = websockets.NewGameSocket(gameUsecase.Queries.FindGameQueryHandler)
	router := mux.NewRouter()

	registerApi(router, gameUsecase)
	registerWS(router, gameSocket)
	registerAsset(router)
	printEndpoints(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func registerAsset(router *mux.Router) {
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./resources")))
}

func registerWS(router *mux.Router, ws http.Handler) {
	router.Handle("/live", ws)
}

// registerApi registers the API endpoints.
func registerApi(router *mux.Router, cases game.UseCases) {
	// main api endpoint
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middleware.HandleError, middleware.JSONContent, middleware.Logger)

	gameRouteBuilder := &RouteBuilder{router: apiRouter}
	gameRouteBuilder.MakeRoute("/game", func(routeBuilder RouteBuilder, router *mux.Router) {
		router.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
			request := game.CreateGameRequestDto{}
			jsonError := json.NewDecoder(r.Body).Decode(&request)
			if jsonError != nil {
				fmt.Printf("%+v\n", jsonError)
			}
			// TODO: Store the settings
		})

		router.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
			request := game.JoinOrLeaveGameRequestDto{}
			jsonError := json.NewDecoder(r.Body).Decode(&request)
			if jsonError != nil {
				fmt.Printf("%+v\n", jsonError)
			}
			// TODO: Store the settings
		})

		router.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
			request := game.StartGameRequestDto{}
			jsonError := json.NewDecoder(r.Body).Decode(&request)
			if jsonError != nil {
				fmt.Printf("%+v\n", jsonError)
			}
			// TODO: Generate words to be guessed
			// TODO: Initialize player sessions
		})

		router.HandleFunc("/leave", func(w http.ResponseWriter, r *http.Request) {
			request := game.JoinOrLeaveGameRequestDto{}
			jsonError := json.NewDecoder(r.Body).Decode(&request)
			if jsonError != nil {
				fmt.Printf("%+v\n", jsonError)
			}
			// TODO: Mark session as destroyed with a reason
		})

		router.HandleFunc("/{id: [0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			gameId := mux.Vars(r)["id"]
			cases.Queries.FindGameQueryHandler.Handle(game.FindGameQuery{ID: uuid.Must(uuid.Parse(gameId))})
		})
	})

	// words endpoint
	//wordsRouter := apiRouter.PathPrefix("/words").Subrouter()

	// auth endpoints
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	auth.RegisterHTTPHandler(authRouter)
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

// loadConfig loads the config from the environment variables or vault
func loadConfig() (*Config, error) {
	conf := &Config{
		DB: NewDBConfig(),
	}
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(conf.DB)
	return conf, err
}

//

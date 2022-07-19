// Package app represents the entry point of the application.
// This is where dependencies are wired together and the application is started.
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lordvidex/wordle-wf/internal/api/handlers"
	"github.com/lordvidex/wordle-wf/internal/auth"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/wordle-wf/internal/adapters"
	"github.com/lordvidex/wordle-wf/internal/api"
	"github.com/lordvidex/wordle-wf/internal/db/pg"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/lordvidex/wordle-wf/internal/middleware"
	"github.com/lordvidex/wordle-wf/internal/websockets"
	"github.com/lordvidex/wordle-wf/internal/words"
	"github.com/spf13/viper"
)

func connectDB(c *DBConfig) (*pgx.Conn, error) {
	var dsn string
	if c.Url != "" {
		dsn = c.Url
	} else {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, 5432, c.DBName)
	}
	pgConn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
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
	authRepo := pg.NewUserRepository(pgConn)

	// services and dependencies
	var gameSocket *websockets.GameSocket
	defer func() {
		if err := gameSocket.Close(); err != nil {
			log.Println("error closing websocket", err)
		}
	}()

	// usecases and application layer components
	wordsUsecase := words.NewUseCases(adapters.NewLocalStringGenerator())
	authUsecase := auth.NewUseCases(
		authRepo,
		adapters.NewPASETOTokenHelper(conf.Token.PASETOSecret, time.Hour),
		adapters.NewBcryptHelper(),
	)
	gameUsecase := game.NewUseCases(
		gameRepo,
		wordsUsecase.RandomWordHandler,
		adapters.NewUniUriGenerator(),
		adapters.NewAwardSystem(),
		gameSocket,
	)

	// adapters and external services
	gameSocket = websockets.NewGameSocket(gameUsecase.Queries.FindGameQueryHandler)
	router := mux.NewRouter()

	registerApi(router, gameUsecase, authUsecase)
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

func routerGroup(parent *mux.Router, path string) *mux.Router {
	return parent.PathPrefix(path).Subrouter()
}

// registerApi registers the API endpoints.
func registerApi(router *mux.Router, gameCases game.UseCases, authCases auth.UseCases) {
	// middlewares
	authMiddleware := api.AuthMiddleware(authCases.Queries.GetUserByToken)
	// main api endpoint
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(
		middleware.HandleError,
		middleware.JSONContent,
		middleware.Logger,
	)

	gameRouter := routerGroup(apiRouter, "/game")
	gameRouter.Use(authMiddleware)
	grh := handlers.NewGameHandler(gameCases)
	gameRouter.HandleFunc("/lobby", grh.CreateLobbyHandler).Methods("POST")
	// TODO(@Israel) - id will be string because of UUID and there will be route clash between this and /lobby
	gameRouter.HandleFunc("/{id: [0-9]+}", grh.GetGameHandler).Methods("GET")

	// auth endpoints
	authRouter := routerGroup(apiRouter, "/auth")
	ah := handlers.NewAuthHandler(authCases)
	authRouter.HandleFunc("/register", ah.RegisterHandler).Methods("POST")
	authRouter.HandleFunc("/login", ah.LoginHandler).Methods("POST")
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

// loadConfig reads the environment variables into *Config
func loadConfig() (*Config, error) {
	conf := NewConfig()
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(conf.DB)
	err = viper.Unmarshal(conf.Token)
	return conf, err
}

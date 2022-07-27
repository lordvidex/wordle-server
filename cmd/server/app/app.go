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
	"github.com/lordvidex/wordle-wf/internal/api"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/lordvidex/wordle-wf/internal/websockets"
	"github.com/sirupsen/logrus"
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
	conf, err := NewConfig()
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

	// services and dependencies
	gameSocket := websockets.NewGameSocket(nil)
	defer func() {
		if err = gameSocket.Close(); err != nil {
			log.Println("error closing websocket", err)
		}
	}()
	wordsUseCase := injectWord()
	authUsecase := injectAuth(pgConn, conf.Token.PASETOSecret, time.Hour)
	gameUsecase := injectGame(pgConn, wordsUseCase.RandomWordHandler, gameSocket)
	gameSocket.
	router := mux.NewRouter()

	registerAPIEndpoints(router, gameUsecase, authUsecase)
	registerWebSocketHandler(router, gameSocket, authUsecase)
	registerAsset(router)
	printAPIEndpoints(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func registerAsset(router *mux.Router) {
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./resources")))
}

func registerWebSocketHandler(router *mux.Router, ws http.Handler, authCases auth.UseCases) {
	authMiddleware := api.AuthMiddleware(authCases.Queries.GetUserByToken)
	router.Path("/live").Handler(authMiddleware(ws))
}

func routerGroup(parent *mux.Router, path string) *mux.Router {
	return parent.PathPrefix(path).Subrouter()
}

func registerAPIEndpoints(router *mux.Router, gameCases game.UseCases, authCases auth.UseCases) {
	// middlewares
	authMiddleware := api.AuthMiddleware(authCases.Queries.GetUserByToken)

	// main api endpoint
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(
		api.HandleError,
		api.JSONContent,
		api.Logger,
	)

	gameRouter := routerGroup(apiRouter, "/game")
	gameRouter.Use(authMiddleware)
	grh := handlers.NewGameHandler(gameCases)
	gameRouter.HandleFunc("/lobby", grh.CreateLobbyHandler).Methods("POST")
	gameRouter.HandleFunc("/{id}", grh.GetGameHandler).Methods("GET")

	// auth endpoints
	authRouter := routerGroup(apiRouter, "/auth")
	ah := handlers.NewAuthHandler(authCases)
	authRouter.HandleFunc("/register", ah.RegisterHandler).Methods("POST")
	authRouter.HandleFunc("/login", ah.LoginHandler).Methods("POST")
}

func printAPIEndpoints(r *mux.Router) {
	if err := r.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		logrus.WithFields(logrus.Fields{
			"methods": methods,
		}).Info(path)
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}

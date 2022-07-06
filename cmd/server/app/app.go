// app represents the entry point of the application.
// This is where dependencies are wired together and the application is started.
package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lordvidex/wordle-wf/internal/auth"
	"github.com/lordvidex/wordle-wf/internal/middleware"
	"github.com/lordvidex/wordle-wf/internal/words"
	"log"
	"net/http"
)

func Start() {
	router := mux.NewRouter()
	registerApi(router)
	registerWS(router)
	registerAsset(router)
	printEndpoints(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func registerAsset(router *mux.Router) {
	router.Handle("/", http.FileServer(http.Dir("./resources")))
}

func registerWS(router *mux.Router) {
	wordsSocket := words.NewWebsocketHandler()
	router.Handle("/live", wordsSocket)
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

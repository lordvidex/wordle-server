package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lordvidex/wordle-wf/internal/api"
	"github.com/lordvidex/wordle-wf/internal/game"
)

type gameRouter struct {
	gameCases game.UseCases
}

func NewGameHandler(gameCases game.UseCases) *gameRouter {
	return &gameRouter{gameCases}
}
func (g *gameRouter) CreateLobbyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	request := game.CreateLobbyRequestDto{}
	jsonError := json.NewDecoder(r.Body).Decode(&request)
	if jsonError != nil {
		fmt.Printf("%+v\n", jsonError)
	}
	lobbyId, err := g.gameCases.Commands.CreateLobbyHandler.Handle(request)
	if err != nil {
		api.BadRequest(err.Error()).WriteJSON(w)
		return
	}
	json.NewEncoder(w).Encode(lobbyId)
}

func (g *gameRouter) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	gameId := mux.Vars(r)["id"]
	gameUUID, err := uuid.Parse(gameId)
	if err != nil {
		api.BadRequest(err.Error()).WriteJSON(w)
	}
	g.gameCases.Queries.FindGameQueryHandler.Handle(game.FindGameQuery{ID: gameUUID})
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lordvidex/wordle-wf/internal/api"
	"github.com/lordvidex/wordle-wf/internal/game"
	"github.com/sirupsen/logrus"
)

type GameResponse struct {
	ID          uuid.UUID     `json:"id"`
	InviteID    string        `json:"invite_id"`
	Word        string        `json:"word,omitempty"`
	Settings    game.Settings `json:"settings"`
	PlayerCount int           `json:"player_count"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     *time.Time    `json:"end_time,omitempty"`
	CreatorID   uuid.UUID     `json:"creator_id,omitempty"`
}

type gameRouter struct {
	gameCases game.UseCases
}

func NewGameHandler(gameCases game.UseCases) *gameRouter {
	return &gameRouter{gameCases}
}
func (g *gameRouter) CreateLobbyHandler(w http.ResponseWriter, r *http.Request) {
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
	response := &game.CreateLobbyResponseDto{LobbyID: lobbyId, Message: "Lobby created successfully"}
	json.NewEncoder(w).Encode(response)
}

func (g *gameRouter) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	gameId := mux.Vars(r)["id"]
	gameUUID, err := uuid.Parse(gameId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		api.BadRequest(err.Error()).WriteJSON(w)
		return
	}
	gm, err := g.gameCases.Queries.FindGameQueryHandler.Handle(game.FindGameQuery{ID: gameUUID})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		api.BadRequest(fmt.Sprintf("an error occured finding game with id %s", gameUUID)).WriteJSON(w)
		logrus.WithFields(logrus.Fields{"error": err.Error()}).Error("error finding game")
		return
	}
	response := &GameResponse{
		ID:       gm.ID,
		InviteID: gm.InviteID,
		Word: func() string {
			if gm.EndTime == nil {
				return ""
			} else {
				return gm.Word.String()
			}
		}(),
		Settings:    gm.Settings,
		PlayerCount: gm.PlayerCount,
		StartTime:   gm.StartTime,
		EndTime:     gm.EndTime,
		CreatorID:   gm.CreatorID,
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

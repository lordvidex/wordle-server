package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lordvidex/wordle-wf/internal/api"
	"github.com/lordvidex/wordle-wf/internal/auth"
)

type PlayerDto struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Points    int64  `json:"points"`
	IsDeleted bool   `json:"is_deleted"`
}
type LoginRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDto struct {
	Data  PlayerDto `json:"data"`
	Token string    `json:"token"`
}

type RegisterDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authHandler struct {
	authCases auth.UseCases
}

func NewAuthHandler(ac auth.UseCases) *authHandler {
	return &authHandler{ac}
}

func (h *authHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var dto LoginRequestDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	// TODO: validate dto
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		api.BadRequest(err.Error()).WriteJSON(w)
		return
	}
	command := auth.LoginCommand{
		Email:    dto.Email,
		Password: dto.Password,
	}
	result, err := h.authCases.Commands.Login.Handle(command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.InternalServerError(err.Error()).WriteJSON(w)
		return
	}

	response := LoginResponseDto{
		Data: PlayerDto{
			ID:        result.Player.ID.String(),
			Name:      result.Player.Name,
			Email:     result.Player.Email,
			Points:    result.Player.Points,
			IsDeleted: result.Player.IsDeleted,
		},
		Token: string(result.Token),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *authHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

}

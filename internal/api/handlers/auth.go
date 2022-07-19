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

type AuthResponseDto struct {
	Data  PlayerDto `json:"data"`
	Token string    `json:"token"`
}
func (a *AuthResponseDto) ParseAuthResponse(payload auth.PlayerWithToken) {
	a.Data = PlayerDto{
		ID:        payload.Player.ID.String(),
		Name:      payload.Player.Name,
		Email:     payload.Player.Email,
		Points:    payload.Player.Points,
		IsDeleted: payload.Player.IsDeleted,
	}
	a.Token = string(payload.Token)
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
		if err == auth.ErrInvalidPassword || err == auth.ErrUserNotFound {
			w.WriteHeader(http.StatusBadRequest)
			api.BadRequest(err.Error()).WriteJSON(w)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			api.InternalServerError(err.Error()).WriteJSON(w)
		}
		return
	}

	var response AuthResponseDto
	response.ParseAuthResponse(*result)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *authHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var dto RegisterDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		api.BadRequest(err.Error()).WriteJSON(w)
		return
	}
	command := auth.RegisterCommand{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
	}
	result, err := h.authCases.Commands.Register.Handle(command)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		api.InternalServerError(err.Error()).WriteJSON(w)
		return
	}
	var response AuthResponseDto
	response.ParseAuthResponse(*result)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

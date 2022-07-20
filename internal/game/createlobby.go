package game

import "fmt"

type CreateLobbyRequestDto struct {
	WordLength            NullInt64 `json:"wordLength"`
	Trials                NullInt64 `json:"trials"`
	PlayerCount           NullInt64 `json:"maxPlayers"`
	Analytics             bool      `json:"hasAnalytics"`
	RecordTime            bool      `json:"shouldRecordTime"`
	ViewOpponentsSessions bool      `json:"canViewOpponentsSessions"`
}

type CreateLobbyHandler interface {
	Handle(lobby CreateLobbyRequestDto) (string, error)
}

type createLobbyHandler struct {
	inviteIdGenerator InviteIDGenerator
	lobbyCreator      LobbyCreator
}

func NewCreateLobbyHandler(i InviteIDGenerator, l LobbyCreator) CreateLobbyHandler {
	return &createLobbyHandler{i, l}
}

func (h *createLobbyHandler) Handle(lobby CreateLobbyRequestDto) (string, error) {
	// validate player count
	if lobby.PlayerCount.Valid && lobby.PlayerCount.Int64 > 10 {
		return "", fmt.Errorf("players in a lobby cannot be more than 10")
	}

	//initialize default values
	if !lobby.PlayerCount.Valid {
		_ = lobby.PlayerCount.Scan(10)
	}
	if !lobby.WordLength.Valid {
		_ = lobby.WordLength.Scan(5)
	}
	if !lobby.Trials.Valid {
		_ = lobby.Trials.Scan(6)
	}

	settings := &Settings{
		MaxPlayerCount:        int(lobby.PlayerCount.Int64),
		WordLength:            int(lobby.WordLength.Int64),
		Trials:                int(lobby.Trials.Int64),
		Analytics:             lobby.Analytics,
		RecordTime:            lobby.RecordTime,
		ViewOpponentsSessions: lobby.ViewOpponentsSessions,
	}
	return h.lobbyCreator.CreateLobby(settings, h.inviteIdGenerator.Generate())
}

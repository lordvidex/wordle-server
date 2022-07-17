package game

import (
	"github.com/lordvidex/wordle-wf/internal/words"
)

type Queries struct {
	FindGameQueryHandler FindGameByIDQueryHandler
}

type Commands struct {
	CreateGameHandler CreateGameHandler
	SubmitResultHandler SubmitResultHandler
	CreateLobbyHandler CreateLobbyHandler
}

type UseCases struct {
	Queries  Queries
	Commands Commands
}

func NewUseCases(repo Repository, g words.RandomHandler,i InviteIDGenerator, a AwardSystem, l LobbyCreator) UseCases {
	return UseCases{
		Queries: Queries{
			FindGameQueryHandler: NewFindGameByIDQueryHandler(repo),
		},
		Commands: Commands{
			SubmitResultHandler: NewSubmitResultHandler(repo, a),
			CreateGameHandler: NewCreateGameHandler(repo, g),
			CreateLobbyHandler: NewCreateLobbyHandler(i, l),
		},
	}
}

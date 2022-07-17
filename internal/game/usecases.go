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
	StartGameHandler   StartGameHandler
	CreateLobbyHandler CreateLobbyHandler
}

type UseCases struct {
	Queries  Queries
	Commands Commands
}

func NewUseCases(repo Repository, g words.RandomHandler,i InviteIDGenerator, n NotificationService, a AwardSystem) UseCases {
	return UseCases{
		Queries: Queries{
			FindGameQueryHandler: NewFindGameByIDQueryHandler(repo),
		},
		Commands: Commands{
			SubmitResultHandler: NewSubmitResultHandler(repo, a),
			CreateGameHandler: NewCreateGameHandler(repo),
			StartGameHandler:  NewStartGameCommandHandler(repo, g, n),
			CreateLobbyHandler: NewCreateLobbyHandler(i),
		},
	}
}

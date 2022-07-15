package game

import (
	"github.com/lordvidex/wordle-wf/internal/words"
)

type Queries struct {
	FindGameQueryHandler     FindGameByIDQueryHandler
	FindGameByInviteHandler  FindByInviteIDQueryHandler
	FindAllGamesQueryHandler FindAllGamesQueryHandler
}

type Commands struct {
	StartGameHandler  StartGameHandler
	CreateGameHandler CreateGameHandler
}

type UseCases struct {
	Queries  Queries
	Commands Commands
}

func NewUseCases(repo Repository, g words.RandomHandler, i InviteIDGenerator, n NotificationService) UseCases {
	return UseCases{
		Queries: Queries{
			FindGameQueryHandler:     NewFindGameByIDQueryHandler(repo),
			FindGameByInviteHandler:  NewFindByInviteIDQueryHandler(repo),
			FindAllGamesQueryHandler: NewFindAllGamesQueryHandler(repo),
		},
		Commands: Commands{
			StartGameHandler:  NewStartGameCommandHandler(repo, g, n),
			CreateGameHandler: NewCreateGameHandler(repo, i),
		},
	}
}

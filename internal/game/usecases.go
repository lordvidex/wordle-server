package game

import (
	"github.com/lordvidex/wordle-wf/internal/words"
)

type Queries struct {
	FindGameQueryHandler     FindGameQueryHandler
	FindAllGamesQueryHandler FindAllGamesQueryHandler
}

type Commands struct {
	StartGameCommandHandler StartGameCommandHandler
}

type UseCases struct {
	Queries  Queries
	Commands Commands
}

func NewUseCases(repo Repository, g words.RandomHandler, n NotificationService) UseCases {
	return UseCases{
		Queries: Queries{
			FindGameQueryHandler:     NewFindGameQueryHandler(repo),
			FindAllGamesQueryHandler: NewFindAllGamesQueryHandler(repo),
		},
		Commands: Commands{
			StartGameCommandHandler: NewStartGameCommandHandler(repo, g, n),
		},
	}
}

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
}

type UseCases struct {
	Queries  Queries
	Commands Commands
}

func NewUseCases(repo Repository, g words.RandomHandler, a AwardSystem) UseCases {
	return UseCases{
		Queries: Queries{
			FindGameQueryHandler: NewFindGameByIDQueryHandler(repo),
		},
		Commands: Commands{
			SubmitResultHandler: NewSubmitResultHandler(repo, a),
			CreateGameHandler: NewCreateGameHandler(repo),
		},
	}
}

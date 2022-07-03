//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/lordvidex/wordle-wf/internal/words"
)

func RegisterWordsHandler(parentRouter *mux.Router) *words.Handler {
	wire.Build(words.NewRepository, words.NewService, words.RegisterHandler)
	return nil
}

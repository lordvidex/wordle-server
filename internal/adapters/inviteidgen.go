package adapters

import (
	"github.com/dchest/uniuri"
	"github.com/lordvidex/wordle-wf/internal/game"
)

type uniURIGenerator struct {
}

func (u *uniURIGenerator) Generate() string {
	return uniuri.NewLen(6)
}

func NewUniUriGenerator() game.InviteIDGenerator {
	return &uniURIGenerator{}
}

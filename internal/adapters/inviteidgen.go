package adapters

import (
	"github.com/dchest/uniuri"
	"github.com/lordvidex/wordle-wf/internal/game"
)

type uniUriGenerator struct {
}

func (u *uniUriGenerator) Generate() string {
	return uniuri.NewLen(6)
}

func NewUniUriGenerator() game.InviteIDGenerator {
	return &uniUriGenerator{}
}

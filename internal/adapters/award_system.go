package adapters

import "github.com/lordvidex/wordle-wf/internal/game"
type awardSystem struct {}

// AwardPoints is a function that awards points to a player
// it awards 10 points to the first player and 0 to the rest
func (a *awardSystem) AwardPoints(position int) int {
	if position == 1 {
		return 10
	}
	return 0
}

func NewAwardSystem() game.AwardSystem {
	return &awardSystem{}
}

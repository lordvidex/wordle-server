package game

// AwardSystem defines an interface for implementing the points awarded to players
// given their position
type AwardSystem interface {
	AwardPoints(position int) int
}

// InviteIDGenerator defines an interface for generating invite IDs (slugs)
type InviteIDGenerator interface {
	Generate() string
}

type LobbyCreator interface {
	CreateLobby(settings *Settings, id string) (string, error)
}

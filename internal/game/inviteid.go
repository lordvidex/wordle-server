package game

type InviteIDGenerator interface {
	Generate() string
}

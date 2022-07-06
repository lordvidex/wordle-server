package game

// Repository is the interface for storing game properties and state
type Repository interface {
	Create(game *Game) (*Game, error)
	Save(game *Game) (*Game, error)
	Find(id string) (*Game, error)
	FindAll() ([]Game, error)
	Delete(id string) error
	Close() error
}

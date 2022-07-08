package mongo

// PlayerModel represents a model in mongo database
type PlayerModel struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

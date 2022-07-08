package mongo

// WordModel represents a word.Word model in mongo database
type WordModel struct {
	ID string `bson:"_id"`
}

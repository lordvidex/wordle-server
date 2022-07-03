package words

type Repository interface {
	AddWord(string) (int, error)
	GetWordByID(int) (*Word, error)
	GetWordByValue(string) (*Word, error)
	SaveWord(Word) error
	GetWordCount() int
}

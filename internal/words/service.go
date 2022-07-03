package words

import "math/rand"

type Service interface {
	GetRandomWord() Word
	GetWordWithID(int) (Word, error)
}

// service is a simple service that receives a Repository that contains words.
// the job of the service is to return a random one out of the many words
type service struct {
	repo Repository
}

func (s *service) GetRandomWord() Word {
	count := s.repo.GetWordCount()

	word, err := s.repo.GetWordByID(rand.Intn(count))
	if err != nil {
		panic(err)
	}
	return *word
}

func (s *service) GetWordWithID(id int) (Word, error) {
	word, err := s.repo.GetWordByID(id)
	return *word, err
}

func NewService(repository Repository) Service {
	return &service{repository}
}

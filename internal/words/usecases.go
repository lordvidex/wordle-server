package words

type UseCases struct {
	RandomWordHandler RandomHandler
}

func NewUseCases(g StringGenerator) UseCases {
	return UseCases{NewRandomHandler(g)}
}

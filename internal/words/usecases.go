package words

type Queries struct {
	GetRandomWordHandler RandomHandler
	GetWordHandler       GetWordHandler
}

type Commands struct {
	AddWordHandler AddWordHandler
}

type UseCases struct {
	Queries  Queries
	Commands Commands
}

func NewUseCases(g StringGenerator, repo Repository) UseCases {
	return UseCases{
		Queries{
			NewRandomHandler(g),
			NewGetWordHandler(repo),
		},
		Commands{
			NewAddWordHandler(repo),
		},
	}
}

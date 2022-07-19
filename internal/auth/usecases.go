package auth

type Queries struct {
	GetUserByToken GetUserByTokenQueryHandler
}

type Commands struct {
	Login    LoginHandler
	Register RegisterHandler
}

type UseCases struct {
	Queries  Queries
	Commands Commands
}

func NewUseCases(r Repository, t TokenHelper, p PasswordHelper) UseCases {
	return UseCases{
		Queries: Queries{
			GetUserByToken: NewUserTokenDecoder(t),
		},
		Commands: Commands{
			Login:    NewLoginHandler(r, t, p),
			Register: NewRegisterHandler(r, t, p),
		},
	}
}

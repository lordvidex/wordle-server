package auth

type Queries struct {
	GetUserByToken QueryGetUserByToken
}

type Commands struct {
	Login    LoginHandler
	Register RegisterHandler
}

type UseCases struct {
	Queries  Queries
	Commands Commands
}

func NewUseCases(repo Repository, tokenHelper TokenHelper[User], passwordChecker PasswordChecker) UseCases {
	return UseCases{
		Queries: Queries{
			GetUserByToken: NewUserTokenDecoder(tokenHelper),
		},
		Commands: Commands{
			Login:    NewLoginHandler(repo, tokenHelper, passwordChecker),
			Register: NewRegisterHandler(repo, tokenHelper),
		},
	}
}

package auth

type RegisterCommand struct {
	Email    string
	Password string
}

type RegisterHandler interface {
	Handle(command RegisterCommand) (token Token, err error)
}

type registerHandler struct {
	repo           Repository
	tokenGenerator TokenHelper[User]
}

func (h *registerHandler) Handle(command RegisterCommand) (token Token, err error) {
	user := &User{
		Email:    command.Email,
		Password: command.Password,
	}
	user, err = h.repo.Create(user)
	if err != nil {
		return "", err
	}
	return h.tokenGenerator.Generate(*user), nil
}

func NewRegisterHandler(repo Repository, tokenGenerator TokenHelper[User]) RegisterHandler {
	return &registerHandler{repo, tokenGenerator}
}

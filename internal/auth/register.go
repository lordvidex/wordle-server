package auth

type RegisterCommand struct {
	Email    string
	Name     string
	Password string
}

type RegisterHandler interface {
	Handle(command RegisterCommand) (token Token, err error)
}

type registerHandler struct {
	repo           Repository
	tokenGenerator TokenHelper
}

func (h *registerHandler) Handle(command RegisterCommand) (token Token, err error) {
	user := &User{
		Email:    command.Email,
		Name:     command.Name,
		Password: command.Password,
	}
	user, err = h.repo.Create(user)
	if err != nil {
		return "", err
	}
	return h.tokenGenerator.Generate(user), nil
}

func NewRegisterHandler(repo Repository, tokenGenerator TokenHelper) RegisterHandler {
	return &registerHandler{repo, tokenGenerator}
}

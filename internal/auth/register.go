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
	player, err := h.repo.Create(command.Name, command.Email, command.Password)
	if err != nil {
		return "", err
	}
	return h.tokenGenerator.Generate(player)
}

func NewRegisterHandler(repo Repository, tokenGenerator TokenHelper) RegisterHandler {
	return &registerHandler{repo, tokenGenerator}
}

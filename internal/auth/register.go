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
	passwordHelper PasswordHelper
}

func (h *registerHandler) Handle(command RegisterCommand) (token Token, err error) {
	hashedPassword, err := h.passwordHelper.Hash(command.Password)
	if err != nil {
		return "", err
	}
	player, err := h.repo.Create(command.Name, command.Email, hashedPassword)
	if err != nil {
		return "", err
	}
	return h.tokenGenerator.Generate(player)
}

func NewRegisterHandler(repo Repository, tokenGenerator TokenHelper, passwordHelper PasswordHelper) RegisterHandler {
	return &registerHandler{repo, tokenGenerator, passwordHelper}
}

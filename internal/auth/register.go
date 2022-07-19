package auth

type RegisterCommand struct {
	Email    string
	Name     string
	Password string
}

type RegisterHandler interface {
	Handle(command RegisterCommand) (*PlayerWithToken, error)
}

type registerHandler struct {
	repo           Repository
	tokenGenerator TokenHelper
	passwordHelper PasswordHelper
}

func (h *registerHandler) Handle(command RegisterCommand) (*PlayerWithToken, error) {
	hashedPassword, err := h.passwordHelper.Hash(command.Password)
	if err != nil {
		return nil, err
	}
	player, err := h.repo.Create(command.Name, command.Email, hashedPassword)
	if err != nil {
		return nil, err
	}
	token, err := h.tokenGenerator.Generate(player)
	if err != nil {
		return nil, err
	}
	return &PlayerWithToken{player, token}, nil
}

func NewRegisterHandler(repo Repository, tokenGenerator TokenHelper, passwordHelper PasswordHelper) RegisterHandler {
	return &registerHandler{repo, tokenGenerator, passwordHelper}
}

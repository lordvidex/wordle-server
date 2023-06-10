package local

import tea "github.com/charmbracelet/bubbletea"

var _ tea.Model = createAccountModel{}

// createAccountModel is a model just before the start of a game to choose the game configurations
type createAccountModel struct {
}

// Init implements tea.Model.
func (m createAccountModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m createAccountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// handle quit case
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m createAccountModel) View() string {
	return "create account model"
}

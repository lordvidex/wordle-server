package local

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lordvidex/wordle-wf/cmd/client/local/store"
)

var _ tea.Model = introModel{}

var (
	defaultPadding = lipgloss.NewStyle().Padding(1, 3)
)

// introModel is the splashscreen model that is landed upon and checks for authentication of the user
type introModel struct {
	loader        spinner.Model
	width, height int
}

func NewIntroModel() introModel {
	sp := spinner.New()
	sp.Spinner = spinner.Monkey
	return introModel{loader: sp}
}

type userMsg struct {
	user   string
	exists bool
}

// loadUser fetches the user from the store and returns fetchUserResponse
func loadUser() tea.Msg {
	user, exists := store.FetchUser()
	return userMsg{user, exists}
}

// Init implements tea.Model.
func (m introModel) Init() tea.Cmd {
	return tea.Batch(m.loader.Tick, loadUser)
	// TODO:
}

// Update implements tea.Model.
func (m introModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	case userMsg:
		// determine whether the user is an old user or a new user by reading from file
		// if old user show the config screen
		// if new user show create account screen
		if msg.exists {
			return NewConfigModel(m.width, m.height), nil
		} else {
			return createAccountModel{}, nil
		}
	case tea.KeyMsg:
		// handle quit case
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}
	m.loader, _ = m.loader.Update(m.loader.Tick)
	return m, nil
}

// View implements tea.Model.
func (m introModel) View() string {
	str := fmt.Sprintf("Loading... \n, %s", m.loader.View())
	return defaultPadding.Render(str)
}

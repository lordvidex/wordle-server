package local

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = configModel{}

// configModel is a model just before the start of a game to choose the game configurations
type configModel struct {
	ti   textinput.Model
	w, h int
}

func NewConfigModel(width, height int) configModel {
	ti := textinput.New()
	ti.Placeholder = "5"
	ti.CharLimit = 1
	ti.Focus()
	return configModel{ti: ti, w: width, h: height}
}

// Init implements tea.Model.
func (m configModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m configModel) gotoGame(w, h int) (tea.Model, tea.Cmd) {
	if m.ti.Value() == "" {
		return m, nil
	}
	digit, err := strconv.Atoi(m.ti.Value())
	if err != nil {
		log.Println("error converting string to int")
		return m, nil
	}
	game := NewGameModel(digit, w, h)
	return game, nil
}

// Update implements tea.Model.
func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// handle quit case
		// handle enter press
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		} else if msg.Type == tea.KeyEnter {
			return m.gotoGame(m.w, m.h)
		} else if key.Matches(msg, m.ti.KeyMap.DeleteCharacterBackward) {
			m.ti, cmd = m.ti.Update(msg)
			return m, cmd
		}

		// handle invalid cases
		if len(msg.Runes) != 1 {
			return m, nil
		}
		rn := msg.Runes[0]
		if rn <= '0' || rn > '9' {
			return m, nil
		}

	}
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m configModel) View() string {
	return fmt.Sprintf(`Welcome to WORDLE Local!
	How many tries do you want per round?

	%s`, m.ti.View())
}

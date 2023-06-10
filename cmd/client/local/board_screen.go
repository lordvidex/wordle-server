package local

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lordvidex/wordle-wf/internal/words"
)

var _ tea.Model = boardModel{}

var (
	boxStyle = lipgloss.NewStyle().Width(5).Height(5).
			Border(lipgloss.RoundedBorder()).
			Foreground(lipgloss.Color("#ffffff"))
	greyStyle   = boxStyle.Background(lipgloss.Color("#888888"))
	greenStyle  = boxStyle.Background(lipgloss.Color("#00ff00"))
	yellowStyle = boxStyle.Background(lipgloss.Color("#ffff00"))
)

var (
	wordLen = 5
)

// boardModel displays the messages from the chat of the game
// as well as in-game notifications from each user.
type boardModel struct {
	session      *Session
	currentWord  []rune
	currentIndex int

	height, width int
}

func initBoardModel(session *Session, height, width int) boardModel {
	return boardModel{
		session:     session,
		currentWord: make([]rune, wordLen),
		height:      height,
		width:       width,
	}
}

func (m boardModel) checkWin() tea.Cmd {
	return func() tea.Msg {
		if m.session.IsWon() {
			return gameWon
		}
		if m.session.HasEnded() {
			return gameLost
		}
		return gamePlaying
	}
}

// Init implements tea.Model.
func (boardModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m boardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m boardModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyBackspace:
		if m.currentIndex > 0 {
			m.currentIndex--
		}
	case tea.KeyEnter:
		// validations
		if m.currentIndex != wordLen {
			// TODO: show message
			return m, nil
		}
		word := words.New(string(m.currentWord))
		m.session.WordStatus(word)
		m.session.PlayedWords = append(m.session.PlayedWords, word)
		m.currentIndex = 0
		return m, m.checkWin()
	case tea.KeyRunes:
		if m.currentIndex == wordLen {
			return m, nil
		}
		str := msg.Runes
		if len(str) != 1 {
			return m, nil
		}
		rn := str[0]
		switch {
		case rn >= 'a' && rn <= 'z':
			rn -= 32
			fallthrough
		case rn >= 'A' && rn <= 'Z':
			m.currentWord[m.currentIndex] = rn
			m.currentIndex++
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m boardModel) View() string {
	// draw a grid of boxes with i as m.maxTries and j as 5
	lines := make([]string, m.session.MaxTries)
	for i := 0; i < m.session.MaxTries; i++ {
		var msg string
		if i < len(m.session.PlayedWords) {
			// render played words
			word := m.session.PlayedWords[i]
			runes := word.Runes()
			for j := 0; j < wordLen; j++ {
				switch word.Stats[j] {
				case words.Correct:
					msg += greenStyle.Render(string(runes[i]))
				case words.Exists:
					msg += yellowStyle.Render(string(runes[i]))
				default:
					msg += greyStyle.Render(string(runes[i]))
				}
			}
		} else if i == len(m.session.PlayedWords) {
			// render current word
			for j := 0; j < m.currentIndex; j++ {
				msg += greyStyle.Render(string(m.currentWord[j]))
			}
			for j := m.currentIndex; j < wordLen; j++ {
				msg += greyStyle.Render(" ")
			}
		} else {
			// render empty boxes
			for j := 0; j < wordLen; j++ {
				msg += boxStyle.Render(" ")
			}
		}
		lines[i] = msg
	}
	return lipgloss.NewStyle().Width(m.width / 2).Height(m.height / 2).
		Render(strings.Join(lines, "\n"))
}

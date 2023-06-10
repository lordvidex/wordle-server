package local

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lordvidex/wordle-wf/internal/adapters"
	"github.com/lordvidex/wordle-wf/internal/words"
)

var _ tea.Model = gameModel{}

var (
	borderStyle       = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Margin(2)
	activeBorderStyle = borderStyle.Foreground(lipgloss.Color("#00ff00"))
)

// sessionState indicates the active tab on the game page
type sessionState int

const (
	board sessionState = iota
	chat
)

// gameMsg indicates the current state of the game
type gameMsg int

const (
	gamePlaying gameMsg = iota
	gameWon
	gameLost
)

// gameModel is the main model for the game screen
// it contains two children, the board screen and the chat screen
type gameModel struct {
	tries int
	state sessionState // whether to focus on the game or on the chat

	width  int
	height int

	board boardModel
	chat  chatModel

	gameState gameMsg
}

// Init implements tea.Model.
func (gameModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case gameMsg:
		// TODO: handle game win or game loss state
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab:
			m.state = (m.state + 1) % 2
			return m, nil
		default:
			break
		}
	}
	// propagate to appropriate child
	var child tea.Model
	switch m.state {
	case board:
		child, cmd = m.board.Update(msg)
		m.board = child.(boardModel)
	case chat:
		child, _ = m.chat.Update(msg)
		m.chat = child.(chatModel)
	}
	return m, cmd
}

// View implements tea.Model.
func (m gameModel) View() string {
	result := ""
	switch m.gameState {
	case gameWon:
		result = "Yay!! You won!!"
	case gameLost:
		result = "Sorry, you lost :("
	}
	var boardMsg = m.board.View()
	var chatMsg = m.chat.View()
	if m.state == board {
		boardMsg = activeBorderStyle.Render(boardMsg)
		chatMsg = borderStyle.Render(chatMsg)
	} else {
		chatMsg = activeBorderStyle.Render(chatMsg)
		boardMsg = borderStyle.Render(boardMsg)

	}
	return fmt.Sprintf("%s%s\n\n%s", boardMsg, chatMsg, result)
}

func NewGameModel(tries, width, height int) gameModel {
	stringGenerator := adapters.NewLocalStringGenerator()
	generator := words.NewRandomHandler(stringGenerator)
	session := NewSession(tries, generator)
	return gameModel{
		tries:  tries,
		width:  width,
		height: height,
		board:  initBoardModel(session, height, width),
		chat:   initChatModel(width, height),
	}
}

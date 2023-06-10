package local

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = chatModel{}

var (
	builder strings.Builder
)

var (
	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	// serverStyle should be in grey and smaller letter
	serverStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#ff00ff")).Italic(true)
	receiverStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
)

// chatModel displays the messages from the chat of the game
// as well as in-game notifications from each user.
type chatModel struct {
	viewport viewport.Model
	chatbox  textarea.Model

	Username string // the username of the current player connected to this instance

	// used to store the messages in the viewport
	messages []string

	width, height int
}

// Init implements tea.Model.
func (chatModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width/2, msg.Height/2
		m.viewport.Height = (m.height * 8) / 10
		m.viewport.Width = m.width
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// send message
			m.addMessage(Message{Username: m.Username, Content: m.chatbox.Value()})
			m.chatbox.SetValue("")
		}
	}
	// update viewport and chatbox
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.chatbox, cmd = m.chatbox.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m chatModel) View() string {
	return lipgloss.NewStyle().Width(m.width / 2).Height(m.height).
		Render(lipgloss.JoinVertical(lipgloss.Top, m.viewport.View(), m.chatbox.View()))
}

func (m chatModel) addMessage(message Message) {
	var msg string
	switch message.Username {
	case m.Username:
		msg = senderStyle.Render(message.Content)
	case "":
		msg = serverStyle.Render(message.Content)
	default:
		msg = receiverStyle.Render(message.Content)
	}
	m.messages = append(m.messages, msg)
	m.viewport.SetContent(msg)
	m.viewport.GotoBottom()
}

func initChatModel(width, height int) chatModel {
	// bottom text area
	txtBox := textarea.New()
	txtBox.Placeholder = "Send a message..."
	txtBox.Focus()
	txtBox.SetWidth(width / 2)
	txtBox.CharLimit = 280
	txtBox.ShowLineNumbers = false

	// chat model
	m := chatModel{
		viewport: viewport.New(width/2, height), // message area
		chatbox:  txtBox,
		width:    width,
		height:   height,
	}
	msgs := []Message{
		{Content: "Welcome to the wordle chat room"},
		{Username: "Martins", Content: "I will play first"},
		{Username: "lordvidex", Content: "Cool bro"},
		{Content: "Game has started"},
	}
	// default messages
	for _, msg := range msgs {
		m.addMessage(msg)
	}
	return m
}

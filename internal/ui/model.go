package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateChat state = iota
	stateApprovalModal
)

// Model is the main TUI model.
type Model struct {
	PaletteName string
	Styles      Styles
	// Left: conversation
	Messages []string
	// Right: plan / current agent
	PlanText string
	AgentName string
	// Bottom: input
	Input textinput.Model
	// Modal
	State       state
	ApprovalMsg string
	ApprovalAgent string
	ApprovalTool string
	ApprovalArgs string
	// Size
	Width  int
	Height int
}

// NewModel returns an initial model.
func NewModel(paletteName string) Model {
	ti := textinput.New()
	ti.Placeholder = "Type a message..."
	ti.CharLimit = 4096
	return Model{
		PaletteName: paletteName,
		Styles:      NewStyles(paletteName),
		Messages:    []string{},
		Input:       ti,
		State:       stateChat,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type msgAppend string
type msgApprove string

func MsgAppend(s string) tea.Msg { return msgAppend(s) }
func MsgApprove(scope string) tea.Msg { return msgApprove(scope) }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.State == stateApprovalModal {
			switch msg.String() {
			case "d", "esc", "q":
				return m, func() tea.Msg { return MsgApprove("deny") }
			case "o":
				return m, func() tea.Msg { return MsgApprove("once") }
			case "s":
				return m, func() tea.Msg { return MsgApprove("session") }
			case "a":
				return m, func() tea.Msg { return MsgApprove("repo") }
			}
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if m.State == stateChat {
			var cmd tea.Cmd
			m.Input, cmd = m.Input.Update(msg)
			if msg.String() == "enter" {
				text := m.Input.Value()
				m.Input.SetValue("")
				m.Messages = append(m.Messages, "You: "+text)
				return m, tea.Batch(cmd, func() tea.Msg { return msgAppend("") })
			}
			return m, cmd
		}
	case msgAppend:
		// Streamed token or assistant message
		if string(msg) != "" {
			m.Messages = append(m.Messages, string(msg))
		}
		return m, nil
	case msgApprove:
		m.State = stateChat
		m.ApprovalMsg = ""
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if m.Width <= 0 {
		m.Width = 80
	}
	if m.Height <= 0 {
		m.Height = 24
	}
	if m.State == stateApprovalModal {
		return m.viewApprovalModal()
	}
	return m.viewChat()
}

func (m Model) viewChat() string {
	leftW := m.Width/2 - 2
	rightW := m.Width - leftW - 4
	conv := strings.Join(m.Messages, "\n")
	if len(conv) > (m.Height-4)*leftW/2 {
		conv = conv[len(conv)-(m.Height-4)*leftW/2:]
	}
	left := m.Styles.Panel.Width(leftW).Height(m.Height - 4).Render("Chat\n" + conv)
	rightContent := "Plan\n" + m.PlanText
	if m.AgentName != "" {
		rightContent += "\n\nAgent: " + m.AgentName
	}
	right := m.Styles.Panel.Width(rightW).Height(m.Height - 4).Render(rightContent)
	row := lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right)
	inputRow := m.Styles.Border.Render("Input: ") + m.Input.View()
	return lipgloss.JoinVertical(lipgloss.Left, row, inputRow)
}

func (m Model) viewApprovalModal() string {
	body := fmt.Sprintf("[Approval] agent=%s tool=%s\nargs=%s\n\n[d]eny [o]nce [s]ession [a]lways repo", m.ApprovalAgent, m.ApprovalTool, m.ApprovalArgs)
	return m.Styles.Panel.Width(60).Render(body)
}

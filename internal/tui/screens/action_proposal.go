package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentleman-programming/argus/internal/components/scanner"
)

type ActionProposalModel struct {
	Action      scanner.ProposedAction
	Choice      bool // true = Authorize, false = Reject
	answered    bool
	width       int
	height      int
	help        help.Model
}

type ActionKeyMap struct {
	Yes   key.Binding
	No    key.Binding
	Quit  key.Binding
}

var Keys = ActionKeyMap{
	Yes: key.NewBinding(
		key.WithKeys("y", "enter"),
		key.WithHelp("y/enter", "authorize"),
	),
	No: key.NewBinding(
		key.WithKeys("n", "esc"),
		key.WithHelp("n/esc", "reject"),
	),
}

func (k ActionKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Yes, k.No}
}

func (k ActionKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Yes, k.No}}
}

func NewActionProposal(action scanner.ProposedAction) *ActionProposalModel {
	return &ActionProposalModel{
		Action: action,
		help:   help.New(),
	}
}

func (m *ActionProposalModel) Init() tea.Cmd {
	return nil
}

func (m *ActionProposalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Yes):
			m.Choice = true
			m.answered = true
			return m, nil
		case key.Matches(msg, Keys.No):
			m.Choice = false
			m.answered = true
			return m, nil
		}
	}
	return m, nil
}

func (m *ActionProposalModel) View() string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF5F00")).
		MarginBottom(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5F5FFF")).
		Padding(1, 2).
		Width(m.width - 10)

	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#AFAFAF"))
	cmdStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#121212")).
		Foreground(lipgloss.Color("#5FFF5F")).
		Padding(0, 1)

	sb.WriteString(titleStyle.Render("🧠 HITL - Action Authorization Required"))
	sb.WriteString("\n")

	content := fmt.Sprintf(
		"Agent wants to perform a %s action:\n\n%s\n\n%s\n%s",
		strings.ToUpper(string(m.Action.Type)),
		descStyle.Render("Description: "+m.Action.Description),
		descStyle.Render("Action:"),
		cmdStyle.Render(m.Action.Command),
	)
	
	if m.Action.Type == scanner.ActionWrite {
		content = fmt.Sprintf(
			"Agent wants to perform a %s action:\n\n%s\n\n%s: %s\n\n%s\n%s",
			strings.ToUpper(string(m.Action.Type)),
			descStyle.Render("Description: "+m.Action.Description),
			descStyle.Render("Path"), m.Action.Path,
			descStyle.Render("Content:"),
			cmdStyle.Render("(Content hidden for brevity, see dash)"),
		)
	}

	sb.WriteString(boxStyle.Render(content))
	sb.WriteString("\n\n")
	sb.WriteString(m.help.View(Keys))

	return sb.String()
}

func (m *ActionProposalModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *ActionProposalModel) Answered() bool {
	return m.answered
}

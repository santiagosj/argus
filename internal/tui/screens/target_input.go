package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/tui/styles"
)

type TargetInputModel struct {
	textInput textinput.Model
	category  agents.NISTCategory
	done      bool
	width     int
	height    int
}

func NewTargetInput() *TargetInputModel {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 2048
	ti.Width = 80

	return &TargetInputModel{
		textInput: ti,
	}
}

func (m *TargetInputModel) SetCategory(cat agents.NISTCategory) {
	m.category = cat
	if cat == agents.Evolve {
		m.textInput.Placeholder = "What should Argus learn? (e.g. new port scanner skill)"
	} else {
		m.textInput.Placeholder = "192.168.1.1 or example.com"
	}
	m.textInput.Reset()
}

func (m *TargetInputModel) Reset() {
	m.textInput.Reset()
	m.done = false
}

func (m *TargetInputModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m TargetInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TargetInputModel) Done() bool {
	return m.done
}

func (m TargetInputModel) Value() string {
	return m.textInput.Value()
}

func (m *TargetInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.textInput.Value() != "" {
				m.done = true
			}
			return m, nil
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m TargetInputModel) View() string {
	var prompt string
	if m.category == agents.Evolve {
		prompt = "What do you want Argus to learn/develop today?"
	} else {
		prompt = "Enter the target IP or Domain for the NIST workflow:"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		styles.GetBanner(),
		styles.TitleStyle.PaddingLeft(2).Render(fmt.Sprintf("Phase: %s", m.category)),
		"",
		lipgloss.NewStyle().PaddingLeft(2).Render(prompt),
		"",
		lipgloss.NewStyle().PaddingLeft(2).Render(m.textInput.View()),
		"",
		styles.HelpStyle.PaddingLeft(2).Render("Press Enter to confirm • Esc to go back • Ctrl+C to exit"),
	)
}

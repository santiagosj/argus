package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/tui/styles"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	category agents.NISTCategory
	title    string
	desc     string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type NISTSelectorModel struct {
	list     list.Model
	choice   agents.NISTCategory
	quitting bool
	width    int
	height   int
}

func NewNISTSelector() *NISTSelectorModel {
	items := []list.Item{
		item{category: agents.Identify, title: "Identify (ID)", desc: "Inventory assets, vulnerabilities, and risk management."},
		item{category: agents.Protect, title: "Protect (PR)", desc: "Apply safeguards to ensure delivery of critical services."},
		item{category: agents.Detect, title: "Detect (DE)", desc: "Identify the occurrence of a cybersecurity event."},
		item{category: agents.Respond, title: "Respond (RS)", desc: "Take action regarding a detected cybersecurity incident."},
		item{category: agents.Recover, title: "Recover (RC)", desc: "Maintain resilience and restore any capabilities."},
		item{category: agents.Evolve, title: "Evolve (EV)", desc: "Self-development: Generate new agents and skills via SDD flow."},
		item{category: agents.General, title: "General (GN)", desc: "General purpose skills and administrative tools."},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Argus: Select NIST Phase"

	return &NISTSelectorModel{list: l}
}

func (m NISTSelectorModel) Init() tea.Cmd {
	return nil
}

func (m NISTSelectorModel) Choice() agents.NISTCategory {
	return m.choice
}

func (m *NISTSelectorModel) Reset() {
	m.choice = ""
}

func (m *NISTSelectorModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	sw, sh := docStyle.GetFrameSize()
	bannerHeight := lipgloss.Height(styles.GetBanner())
	m.list.SetSize(w-sw, h-sh-bannerHeight)
}

func (m *NISTSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.category
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		h, v := docStyle.GetFrameSize()
		bannerHeight := lipgloss.Height(styles.GetBanner())
		m.list.SetSize(msg.Width-h, msg.Height-v-bannerHeight)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m NISTSelectorModel) View() string {
	if m.choice != "" {
		return fmt.Sprintf("\n  Phase Selected: %s. Launching Autonomous Engine...\n", m.choice)
	}
	if m.quitting {
		return "\n  Exiting Argus. Stay safe.\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		styles.GetBanner(),
		docStyle.Render(m.list.View()),
	)
}

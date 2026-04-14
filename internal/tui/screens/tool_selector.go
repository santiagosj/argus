package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/tui/styles"
)

type toolItem struct {
	name        string
	description string
	category    agents.NISTCategory
}

func (i toolItem) Title() string       { return i.name }
func (i toolItem) Description() string { return i.description }
func (i toolItem) FilterValue() string { return i.name + " " + i.description }

type ToolSelectorModel struct {
	list           list.Model
	target         string
	choice         string
	choiceCategory agents.NISTCategory
	width          int
	height         int
}

func NewToolSelector() *ToolSelectorModel {
	items := []list.Item{
		toolItem{name: "nmap", description: "Network discovery, port/service/OS fingerprinting", category: agents.Identify},
		toolItem{name: "nuclei", description: "Template-based vulnerability detection and scanning", category: agents.Detect},
		toolItem{name: "feroxbuster", description: "Web content and directory discovery using brute-force", category: agents.Identify},
		toolItem{name: "gobuster", description: "Directory/DNS brute-forcing and endpoint enumeration", category: agents.Identify},
		toolItem{name: "dirb", description: "Web content discovery with wordlists", category: agents.Identify},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Argus: Select a tool to run against the target"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return &ToolSelectorModel{
		list: l,
	}
}

func (m *ToolSelectorModel) SetTarget(target string) {
	m.target = target
}

func (m *ToolSelectorModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	h, v := docStyle.GetFrameSize()
	bannerHeight := lipgloss.Height(styles.GetBanner())
	m.list.SetSize(w-h, h-v-bannerHeight)
}

func (m ToolSelectorModel) Init() tea.Cmd {
	return nil
}

func (m ToolSelectorModel) Choice() string {
	return m.choice
}

func (m ToolSelectorModel) Category() agents.NISTCategory {
	return m.choiceCategory
}

func (m *ToolSelectorModel) Reset() {
	m.choice = ""
	m.choiceCategory = ""
}

func (m *ToolSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(toolItem)
			if ok {
				m.choice = i.name
				m.choiceCategory = i.category
				return m, nil
			}
		}
		if msg.String() == "esc" {
			m.choice = "BACK"
			return m, nil
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

func (m ToolSelectorModel) View() string {
	var prompt string
	if m.target != "" {
		prompt = fmt.Sprintf("Target: %s\nChoose a tool to execute:", m.target)
	} else {
		prompt = "Choose a tool to execute against the target:"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		styles.GetBanner(),
		styles.TitleStyle.PaddingLeft(2).Render(prompt),
		"",
		docStyle.Render(m.list.View()),
		"",
		styles.HelpStyle.PaddingLeft(2).Render("Enter to run • Esc to cancel • Ctrl+C to exit"),
	)
}

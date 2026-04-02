package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/tui/styles"
)

type skillItem struct {
	name    string
	path    string
	content string
}

func (i skillItem) Title() string       { return i.name }
func (i skillItem) Description() string { return "Skill file: " + i.path }
func (i skillItem) FilterValue() string { return i.name }

type SkillSelectorModel struct {
	list     list.Model
	category agents.NISTCategory
	choice   string
	quitting bool
	ready    bool
	width    int
	height   int
}

func NewSkillSelector() *SkillSelectorModel {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Argus: Select Skill to Learn/Execute"
	return &SkillSelectorModel{list: l}
}

func (m *SkillSelectorModel) SetCategory(cat agents.NISTCategory) {
	m.category = cat
	m.list.Title = fmt.Sprintf("Argus: Skills for %s", cat)

	items := []list.Item{}
	skillDir := filepath.Join("skills", string(cat))

	files, err := os.ReadDir(skillDir)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
				name := strings.TrimSuffix(f.Name(), ".md")
				items = append(items, skillItem{
					name: name,
					path: filepath.Join(skillDir, f.Name()),
				})
			}
		}
	}

	if len(items) == 0 {
		items = append(items, skillItem{name: "No skills found", path: "N/A"})
	}

	m.list.SetItems(items)
	m.ready = true
	m.choice = ""
}

func (m *SkillSelectorModel) Reset() {
	m.choice = ""
}

func (m *SkillSelectorModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	sw, sh := docStyle.GetFrameSize()
	bannerHeight := lipgloss.Height(styles.GetBanner())
	m.list.SetSize(w-sw, h-sh-bannerHeight)
}

func (m SkillSelectorModel) Init() tea.Cmd {
	return nil
}

func (m SkillSelectorModel) Choice() string {
	return m.choice
}

func (m *SkillSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(skillItem)
			if ok && i.path != "N/A" {
				m.choice = i.name
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

func (m SkillSelectorModel) View() string {
	if m.choice != "" && m.choice != "BACK" {
		return fmt.Sprintf("\n  Skill Selected: %s. Launching...\n", m.choice)
	}
	if m.quitting {
		return "\n  Exiting Argus. Stay safe.\n"
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		styles.GetBanner(),
		docStyle.Render(m.list.View()),
	)
}

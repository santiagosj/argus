package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentleman-programming/argus/internal/tui/styles"
)

type ToolRunningModel struct {
	spinner   spinner.Model
	viewport  viewport.Model
	tool      string
	target    string
	output    strings.Builder
	done      bool
	width     int
	height    int
	startTime time.Time
}

type toolOutputMsg struct {
	output string
}

type toolFinishedMsg struct {
	success bool
}

func NewToolRunning() *ToolRunningModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	vp := viewport.New(0, 0)

	return &ToolRunningModel{
		spinner:   s,
		viewport:  vp,
		startTime: time.Now(),
	}
}

func (m *ToolRunningModel) SetTool(tool, target string) {
	m.tool = tool
	m.target = target
	m.output.Reset()
	m.done = false
	m.startTime = time.Now()
}

func (m *ToolRunningModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w
	m.viewport.Height = h - 6 // Leave space for header and footer
}

func (m ToolRunningModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m ToolRunningModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.done = true
			return m, tea.Quit
		case "esc":
			m.done = true
			return m, nil
		}

	case toolOutputMsg:
		m.output.WriteString(msg.output + "\n")
		m.viewport.SetContent(m.output.String())
		m.viewport.GotoBottom()

	case toolFinishedMsg:
		m.done = true
		if msg.success {
			m.output.WriteString("\n✅ Tool execution completed successfully\n")
		} else {
			m.output.WriteString("\n❌ Tool execution failed\n")
		}
		m.viewport.SetContent(m.output.String())
		m.viewport.GotoBottom()
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m ToolRunningModel) View() string {
	var b strings.Builder

	// Header with tool info
	header := fmt.Sprintf("🔧 Running %s on %s", m.tool, m.target)
	b.WriteString(styles.TitleStyle.Render(header) + "\n")

	// Status line with spinner and elapsed time
	elapsed := time.Since(m.startTime).Truncate(time.Second)
	status := fmt.Sprintf("%s Executing... (%s)", m.spinner.View(), elapsed)
	b.WriteString(styles.InfoStyle.Render(status) + "\n\n")

	// Output viewport
	b.WriteString(m.viewport.View() + "\n")

	// Footer with instructions
	if m.done {
		b.WriteString(styles.HelpStyle.Render("Tool finished • Press ESC to return to tool selection • Q to quit"))
	} else {
		b.WriteString(styles.HelpStyle.Render("Running... • ESC to cancel • Q to quit"))
	}

	return b.String()
}

func (m ToolRunningModel) Done() bool {
	return m.done
}

func (m *ToolRunningModel) ReceiveOutput(output string) {
	m.output.WriteString(output + "\n")
	m.viewport.SetContent(m.output.String())
	m.viewport.GotoBottom()
}

func (m *ToolRunningModel) SetFinished(success bool) {
	m.done = true
	if success {
		m.output.WriteString("\n✅ Tool execution completed successfully\n")
	} else {
		m.output.WriteString("\n❌ Tool execution failed\n")
	}
	m.viewport.SetContent(m.output.String())
	m.viewport.GotoBottom()
}

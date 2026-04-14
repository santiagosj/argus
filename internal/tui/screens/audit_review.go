package screens

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gentleman-programming/argus/internal/system"
)

type AuditReviewModel struct {
	list          list.Model
	viewport      viewport.Model
	searchInput   textinput.Model
	hitlValidator system.HitlValidator
	auditLogger   *system.AuditLogger
	filterPhase   string
	filterStatus  string
	width         int
	height        int
	showDetails   bool
	selectedItem  *AuditItem
}

type AuditItem struct {
	entry system.AuditEntry
}

func (a AuditItem) Title() string {
	timestamp := a.entry.Timestamp.Format("15:04:05")
	source := a.entry.Source
	entryType := a.entry.Type

	var status string
	if a.entry.Status != "" {
		status = fmt.Sprintf(" [%s]", a.entry.Status)
	}

	return fmt.Sprintf("%s %s: %s%s", timestamp, source, entryType, status)
}

func (a AuditItem) Description() string {
	// Show a brief description based on entry type
	switch a.entry.Type {
	case "PROPOSAL":
		if content, ok := a.entry.Content.(map[string]interface{}); ok {
			if action, exists := content["action"]; exists {
				return fmt.Sprintf("Proposed: %v", action)
			}
		}
	case "HITL_REQUEST":
		return "Human approval requested"
	case "HITL_DECISION":
		if a.entry.HitlData != nil {
			return fmt.Sprintf("Decision by %s: %s", a.entry.HitlData.ApprovedBy, a.entry.HitlData.Reason)
		}
	case "ACTION":
		if content, ok := a.entry.Content.(map[string]interface{}); ok {
			if tool, exists := content["tool"]; exists {
				return fmt.Sprintf("Tool: %v", tool)
			}
		}
	}
	return "Audit entry"
}

func (a AuditItem) FilterValue() string {
	return a.Title() + " " + a.Description()
}

func NewAuditReviewModel(hitlValidator system.HitlValidator, auditLogger *system.AuditLogger) AuditReviewModel {
	// Initialize search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search audit entries..."
	searchInput.CharLimit = 100

	// Initialize list
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	l := list.New([]list.Item{}, delegate, 0, 0)
	l.Title = "Argus Audit Review"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false) // We'll handle filtering ourselves

	return AuditReviewModel{
		list:          l,
		searchInput:   searchInput,
		hitlValidator: hitlValidator,
		auditLogger:   auditLogger,
		filterPhase:   "ALL",
		filterStatus:  "ALL",
		showDetails:   false,
	}
}

func (m AuditReviewModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadAuditEntries(),
		textinput.Blink,
	)
}

func (m AuditReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-8)
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.showDetails && m.selectedItem != nil {
				// Handle approval/rejection in details view
				return m, nil
			} else {
				// Select item for details
				if i := m.list.Index(); i >= 0 && i < len(m.list.Items()) {
					if item, ok := m.list.Items()[i].(AuditItem); ok {
						m.selectedItem = &item
						m.showDetails = true
						m.updateViewport()
					}
				}
			}
		case "esc":
			if m.showDetails {
				m.showDetails = false
				m.selectedItem = nil
			} else {
				return m, tea.Quit
			}
		case "/":
			m.searchInput.Focus()
		case "tab":
			// Cycle through filters
			m.cycleFilter()
		}

	case auditEntriesLoadedMsg:
		items := make([]list.Item, len(msg.entries))
		for i, entry := range msg.entries {
			items[i] = AuditItem{entry: entry}
		}
		cmd = m.list.SetItems(items)
	}

	// Update components
	if m.searchInput.Focused() {
		m.searchInput, cmd = m.searchInput.Update(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m AuditReviewModel) View() string {
	if m.showDetails && m.selectedItem != nil {
		return m.detailsView()
	}

	var b strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00ff00")).
		Render("🔍 ARGUS AUDIT REVIEW")
	b.WriteString(header + "\n\n")

	// Filters
	filters := fmt.Sprintf("Phase: %s | Status: %s | Search: %s",
		m.filterPhase, m.filterStatus, m.searchInput.View())
	b.WriteString(filters + "\n\n")

	// Instructions
	instructions := "↑/↓ Navigate • Enter Details • / Search • Tab Cycle Filters • Esc/Quit"
	b.WriteString(lipgloss.NewStyle().Faint(true).Render(instructions) + "\n\n")

	// List
	b.WriteString(m.list.View())

	return b.String()
}

func (m AuditReviewModel) detailsView() string {
	if m.selectedItem == nil {
		return "No item selected"
	}

	entry := m.selectedItem.entry
	var b strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00d4ff")).
		Render("📋 AUDIT ENTRY DETAILS")
	b.WriteString(header + "\n\n")

	// Basic info
	b.WriteString(fmt.Sprintf("Timestamp: %s\n", entry.Timestamp.Format(time.RFC3339)))
	b.WriteString(fmt.Sprintf("Source: %s\n", entry.Source))
	b.WriteString(fmt.Sprintf("Type: %s\n", entry.Type))
	b.WriteString(fmt.Sprintf("Status: %s\n", entry.Status))
	b.WriteString("\n")

	// Content
	b.WriteString("Content:\n")
	contentStr := fmt.Sprintf("%+v", entry.Content)
	if len(contentStr) > 200 {
		contentStr = contentStr[:200] + "..."
	}
	b.WriteString(contentStr + "\n\n")

	// HITL Data
	if entry.HitlData != nil {
		b.WriteString("HITL Decision:\n")
		b.WriteString(fmt.Sprintf("  Approved By: %s\n", entry.HitlData.ApprovedBy))
		b.WriteString(fmt.Sprintf("  At: %s\n", entry.HitlData.ApprovedAt.Format(time.RFC3339)))
		b.WriteString(fmt.Sprintf("  Reason: %s\n", entry.HitlData.Reason))
		if entry.HitlData.ActualOutcome != "" {
			b.WriteString(fmt.Sprintf("  Outcome: %s\n", entry.HitlData.ActualOutcome))
		}
		b.WriteString("\n")
	}

	// Actions for pending items
	if entry.Status == "PENDING" && entry.Type == "PROPOSAL" {
		b.WriteString("Actions:\n")
		b.WriteString("  [A]pprove  [R]eject  [Esc] Back\n")
	}

	// Instructions
	b.WriteString(lipgloss.NewStyle().Faint(true).Render("Esc Back • Q Quit"))

	return b.String()
}

func (m *AuditReviewModel) updateViewport() {
	if m.selectedItem != nil {
		content := m.detailsView()
		m.viewport.SetContent(content)
	}
}

func (m *AuditReviewModel) cycleFilter() {
	// Cycle through phase filters
	phases := []string{"ALL", "IDENTIFY", "PROTECT", "DETECT", "RESPOND", "RECOVER", "EVOLVE"}
	currentIndex := 0
	for i, phase := range phases {
		if phase == m.filterPhase {
			currentIndex = i
			break
		}
	}
	m.filterPhase = phases[(currentIndex+1)%len(phases)]

	// Reload with new filter
	m.loadAuditEntries()
}

func (m AuditReviewModel) loadAuditEntries() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual audit loading with filtering
		// For now, return empty list
		return auditEntriesLoadedMsg{entries: []system.AuditEntry{}}
	}
}

type auditEntriesLoadedMsg struct {
	entries []system.AuditEntry
}

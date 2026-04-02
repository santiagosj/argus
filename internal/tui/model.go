package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/components/scanner"
	"github.com/gentleman-programming/argus/internal/tui/screens"
)

type Screen int

const (
	ScreenNISTSelector Screen = iota
	ScreenSkillSelector
	ScreenTargetInput
	ScreenRunning
	ScreenActionProposal
)

// Messages for Orchestrator communication
type ActionRequestMsg struct {
	Action scanner.ProposedAction
	Result chan bool
}

type WorkflowFinishedMsg struct {
	Err error
}

type TargetSelectedMsg struct {
	Category agents.NISTCategory
	Target   string
	Skill    string
}

type Model struct {
	currentScreen  Screen
	nistSelector   *screens.NISTSelectorModel
	skillSelector  *screens.SkillSelectorModel
	targetInput    *screens.TargetInputModel
	actionProposal *screens.ActionProposalModel

	// Communication
	TargetSelectedChan chan TargetSelectedMsg
	currentResultChan  chan bool

	// Global State
	SelectedCategory agents.NISTCategory
	SelectedSkill    string
	Target           string
	Quitting         bool
	LearningMode     bool

	// Window Size
	lastWidth  int
	lastHeight int
}

func NewModel() Model {
	nist := screens.NewNISTSelector()
	skill := screens.NewSkillSelector()
	target := screens.NewTargetInput()

	return Model{
		currentScreen: ScreenNISTSelector,
		nistSelector:  nist,
		skillSelector: skill,
		targetInput:   target,
		lastWidth:     80, // Defaults
		lastHeight:    24,
	}
}

func NewLearningModel() Model {
	m := NewModel()
	m.LearningMode = true
	return m
}

func (m Model) Init() tea.Cmd {
	return m.nistSelector.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.lastWidth = msg.Width
		m.lastHeight = msg.Height
		m.nistSelector.SetSize(msg.Width, msg.Height)
		m.skillSelector.SetSize(msg.Width, msg.Height)
		m.targetInput.SetSize(msg.Width, msg.Height)
		if m.actionProposal != nil {
			m.actionProposal.SetSize(msg.Width, msg.Height)
		}
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	case ActionRequestMsg:
		m.actionProposal = screens.NewActionProposal(msg.Action)
		m.actionProposal.SetSize(m.lastWidth, m.lastHeight)
		m.currentResultChan = msg.Result
		m.currentScreen = ScreenActionProposal
		return m, m.actionProposal.Init()
	case WorkflowFinishedMsg:
		m.currentScreen = ScreenTargetInput // Or a report screen
		return m, nil
	}

	switch m.currentScreen {
	case ScreenNISTSelector:
		newNist, nistCmd := m.nistSelector.Update(msg)
		m.nistSelector = newNist.(*screens.NISTSelectorModel)
		cmd = nistCmd

		if m.nistSelector.Choice() != "" {
			m.SelectedCategory = m.nistSelector.Choice()

			// Si es EVOLVE o estamos en modo aprendizaje, vamos al SkillSelector
			if m.SelectedCategory == agents.Evolve || m.LearningMode {
				m.currentScreen = ScreenSkillSelector
				m.skillSelector.SetCategory(m.SelectedCategory)
				m.skillSelector.SetSize(m.lastWidth, m.lastHeight)
				return m, m.skillSelector.Init()
			}

			m.currentScreen = ScreenTargetInput
			m.targetInput.SetCategory(m.SelectedCategory)
			m.targetInput.Reset()
			return m, m.targetInput.Init()
		}

	case ScreenSkillSelector:
		newSkillSelector, skillCmd := m.skillSelector.Update(msg)
		m.skillSelector = newSkillSelector.(*screens.SkillSelectorModel)
		cmd = skillCmd

		choice := m.skillSelector.Choice()
		if choice == "BACK" {
			m.currentScreen = ScreenNISTSelector
			m.nistSelector.Reset()
			m.nistSelector.SetSize(m.lastWidth, m.lastHeight)
			return m, m.nistSelector.Init()
		} else if choice != "" {
			m.SelectedSkill = choice
			m.currentScreen = ScreenTargetInput
			m.targetInput.SetCategory(m.SelectedCategory)
			m.targetInput.Reset()
			return m, m.targetInput.Init()
		}

	case ScreenTargetInput:
		// Fix: Handle ESC to go back from Target Input
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
			if m.SelectedSkill != "" {
				m.currentScreen = ScreenSkillSelector
				m.skillSelector.Reset()
				return m, m.skillSelector.Init()
			}
			m.currentScreen = ScreenNISTSelector
			m.nistSelector.Reset()
			return m, m.nistSelector.Init()
		}

		newTargetInput, targetCmd := m.targetInput.Update(msg)
		m.targetInput = newTargetInput.(*screens.TargetInputModel)
		cmd = targetCmd

		if m.targetInput.Done() {
			m.Target = m.targetInput.Value()
			m.currentScreen = ScreenRunning
			
			// Notify external orchestrator
			if m.TargetSelectedChan != nil {
				go func() {
					m.TargetSelectedChan <- TargetSelectedMsg{
						Category: m.SelectedCategory,
						Target:   m.Target,
						Skill:    m.SelectedSkill,
					}
				}()
			}
			return m, nil 
		}

	case ScreenActionProposal:
		newProposal, proposalCmd := m.actionProposal.Update(msg)
		m.actionProposal = newProposal.(*screens.ActionProposalModel)
		cmd = proposalCmd

		if m.actionProposal.Answered() {
			m.currentResultChan <- m.actionProposal.Choice
			m.currentScreen = ScreenRunning
			return m, nil
		}
	}

	return m, cmd
}

func (m Model) View() string {
	if m.Quitting {
		return "Exiting Argus...\n"
	}

	switch m.currentScreen {
	case ScreenNISTSelector:
		return m.nistSelector.View()
	case ScreenSkillSelector:
		return m.skillSelector.View()
	case ScreenTargetInput:
		return m.targetInput.View()
	case ScreenRunning:
		return "Launching Autonomous Workflow...\nCheck Dash: http://localhost:8080\n"
	case ScreenActionProposal:
		return m.actionProposal.View()
	default:
		return "Unknown screen"
	}
}

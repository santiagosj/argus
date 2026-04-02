package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/components/engram"
	"github.com/gentleman-programming/argus/internal/components/scanner"
	"github.com/gentleman-programming/argus/internal/system"
	"github.com/gentleman-programming/argus/internal/tui"
	"github.com/gentleman-programming/argus/internal/workflow"
)

var Version = "v3.0.0-alpha"

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		fmt.Printf("Could not open browser: %v\n", err)
	}
}

func Run() error {
	return RunArgs(os.Args[1:])
}

func RunArgs(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "version", "--version", "-v":
			fmt.Printf("Argus %s\n", Version)
			return nil
		case "learn":
			return startTUI(args, true)
		case "help", "--help", "-h":
			fmt.Println("Argus: Cognitive Security Framework")
			fmt.Println("\nUsage:")
			fmt.Println("  argus [command]")
			fmt.Println("\nAvailable Commands:")
			fmt.Println("  run [target]   Start autonomous workflow (opens TUI if no target)")
			fmt.Println("  learn          Open skill learning menu")
			fmt.Println("  version        Show version information")
			return nil
		}
	}

	// Default behavior (TUI)
	return startTUI(args, false)
}

func startTUI(args []string, learningMode bool) error {
	// 1. Initialize Components
	cfg, err := system.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Warning: error loading config: %v. Using defaults.\n", err)
	}

	var provider agents.AgentProvider
	switch cfg.AI.Provider {
	case "ollama":
		provider = agents.NewOllamaProvider(cfg.AI.BaseURL, cfg.AI.Model)
	default:
		// Fallback to Ollama
		provider = agents.NewOllamaProvider(cfg.AI.BaseURL, cfg.AI.Model)
	}

	memory, _ := engram.InjectEngram(context.Background(), cfg.Persistence.Type, cfg.Persistence.Path)
	
	// Start Dashboard
	dashboard := system.NewWebDashboard(8080)
	dashboard.Start()

	// 2. Start TUI
	targetChan := make(chan tui.TargetSelectedMsg)
	var m tui.Model
	if learningMode {
		m = tui.NewLearningModel()
	} else {
		m = tui.NewModel()
	}
	m.TargetSelectedChan = targetChan
	
	p := tea.NewProgram(m, tea.WithAltScreen())

	// Validator function for HITL (Bridge to TUI)
	validate := func(a scanner.ProposedAction) bool {
		resultChan := make(chan bool)
		p.Send(tui.ActionRequestMsg{
			Action: a,
			Result: resultChan,
		})
		return <-resultChan
	}

	orch := workflow.NewOrchestrator(provider, memory, dashboard, validate)

	// Check if target was passed in args (e.g., argus run localhost)
	var initialTarget string
	if len(args) > 1 && args[0] == "run" {
		initialTarget = args[1]
	}

	if initialTarget != "" {
		fmt.Printf("\n--- Launching AUTONOMOUS Workflow on %s ---\n", initialTarget)
		return orch.RunAutonomousWorkflow(context.Background(), initialTarget)
	}

	// Run Orchestrator in a separate goroutine
	go func() {
		// Wait for target selection from TUI
		select {
		case msg := <-targetChan:
			category := msg.Category
			target := msg.Target
			skill := msg.Skill
			
			var err error
			if category == agents.Evolve {
				if skill != "" {
					err = orch.RunEvolutionWorkflow(context.Background(), skill)
				} else {
					err = orch.RunEvolutionWorkflow(context.Background(), target)
				}
			} else {
				err = orch.RunAutonomousWorkflow(context.Background(), target)
			}
			
			p.Send(tui.WorkflowFinishedMsg{Err: err})
		case <-time.After(1 * time.Hour): // Timeout or handle quitting
		}
	}()

	_, err = p.Run()
	return err
}

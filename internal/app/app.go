package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/components/engram"
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
		case "demo":
			return startDemo()
		case "run":
			if len(args) < 2 {
				fmt.Println("Usage: argus run <target>")
				return nil
			}
			return startTUI(args[1:], false)
		case "status":
			return showStatus()
		case "help", "--help", "-h":
			fmt.Println("Argus: Cognitive Security Framework")
			fmt.Println("\nUsage:")
			fmt.Println("  argus                Start the interactive TUI")
			fmt.Println("  argus <ip:port>      Start TUI with target IP/port")
			fmt.Println("  argus learn          Start the interactive learning TUI")
			fmt.Println("  argus demo           Run demo workflow")
			fmt.Println("  argus run <target>   Run workflow on target")
			fmt.Println("\nAvailable Commands:")
			fmt.Println("  learn                Open skill learning menu")
			fmt.Println("  status               Show current status")
			fmt.Println("  version              Show version information")
			fmt.Println("  help                 Show this help message")
			return nil
		default:
			// Check if first arg looks like a target (IP:port or hostname:port)
			if isValidTarget(args[0]) {
				return startTUIWithTarget(args[0], false)
			}
		}
	}

	// Default behavior (TUI)
	return startTUI(args, false)
}

func isValidTarget(target string) bool {
	// Simple validation: contains ":" or is a hostname/IP
	if target == "" {
		return false
	}
	// Allow formats: localhost:3000, 192.168.1.1:8080, example.com:443, etc.
	return len(target) > 0 && !isRunCommand(target)
}

func isRunCommand(s string) bool {
	commands := map[string]bool{
		"version": true, "--version": true, "-v": true,
		"learn": true, "demo": true, "run": true, "status": true,
		"help": true, "--help": true, "-h": true,
	}
	return commands[s]
}

func startWorkflow(target string, learningMode bool) error {
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

	// Create HITL validator for interactive approval
	hitlValidator := system.NewInteractiveValidator()

	orch := workflow.NewOrchestrator(provider, memory, dashboard, hitlValidator)

	fmt.Printf("\n--- Launching Workflow on %s ---\n", target)
	fmt.Println("Dashboard available at: http://localhost:8080")
	fmt.Println("Progress will be logged here...")

	return orch.RunAutonomousWorkflow(context.Background(), target)
}

func showStatus() error {
	cfg, err := system.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		return err
	}

	fmt.Printf("Argus %s\n", Version)
	fmt.Printf("AI Provider: %s (%s)\n", cfg.AI.Provider, cfg.AI.Model)
	fmt.Printf("Persistence: %s at %s\n", cfg.Persistence.Type, cfg.Persistence.Path)
	fmt.Println("Dashboard: http://localhost:8080 (if running)")
	return nil
}

func startDemo() error {
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
		// Fallback to mock provider for demo
		provider = agents.NewOllamaProvider(cfg.AI.BaseURL, cfg.AI.Model)
	}

	memory, _ := engram.InjectEngram(context.Background(), cfg.Persistence.Type, cfg.Persistence.Path)

	// Start Dashboard
	dashboard := system.NewWebDashboard(8080)
	dashboard.Start()

	// Create HITL validator for demo (auto-approve)
	hitlValidator := system.NewInteractiveValidator()

	orch := workflow.NewOrchestrator(provider, memory, dashboard, hitlValidator)

	fmt.Printf("\n--- Lanzando DEMO de Argus ---\n")
	fmt.Println("Dashboard disponible en: http://localhost:8080")

	// Show external IP for WSL/Windows access
	if ip := getExternalIP(); ip != "" {
		fmt.Printf("Desde Windows: http://%s:8080\n", ip)
	}

	fmt.Println("La demo mostrará todas las fases NIST...")

	err = orch.RunDemoWorkflow(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("\n🎉 Demo completada! El dashboard permanece activo en http://localhost:8080")
	fmt.Println("Presiona Ctrl+C para salir...")

	// Mantener el servidor vivo indefinidamente
	select {}
}

func isLikelyTarget(arg string) bool {
	if arg == "" {
		return false
	}
	if strings.Contains(arg, ":") || strings.Contains(arg, ".") {
		return true
	}
	return false
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

	initialTarget := ""
	if len(args) == 1 && isLikelyTarget(args[0]) {
		initialTarget = args[0]
	}

	// 2. Start TUI
	targetChan := make(chan tui.TargetSelectedMsg)
	var m tui.Model
	if learningMode {
		m = tui.NewLearningModel()
	} else {
		m = tui.NewModel()
	}
	m.TargetSelectedChan = targetChan

	if initialTarget != "" {
		m.EnableDirectTargetMode(initialTarget)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	// Communication channels
	toolOutputChan := make(chan string, 100)
	toolFinishedChan := make(chan bool, 1)

	// Create HITL validator for TUI (auto-approves)
	hitlValidator := system.NewTUIValidator()

	orch := workflow.NewOrchestrator(provider, memory, dashboard, hitlValidator)

	// The TUI is the only supported interface for running workflows.
	// Any target selection is handled inside the TUI.

	// Run Orchestrator in a separate goroutine
	go func() {
		for {
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
					// For evolution workflows, finish the TUI
					p.Send(tui.WorkflowFinishedMsg{Err: err})
					return
				} else if skill != "" {
					// For individual tools, run and send updates but don't finish TUI
					err = orch.RunToolOnTargetWithUpdates(context.Background(), target, skill, toolOutputChan)
					toolFinishedChan <- (err == nil)
				} else {
					// For full workflows, finish the TUI
					err = orch.RunAutonomousWorkflow(context.Background(), target)
					p.Send(tui.WorkflowFinishedMsg{Err: err})
					return
				}
			}
		}
	}()

	// Handle real-time updates
	go func() {
		for {
			select {
			case output := <-toolOutputChan:
				p.Send(tui.ToolOutputMsg{Output: output})
			case finished := <-toolFinishedChan:
				p.Send(tui.ToolFinishedMsg{Success: finished})
			}
		}
	}()

	_, err = p.Run()
	return err
}

func getExternalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func startTUIWithTarget(target string, learningMode bool) error {
	return startTUI([]string{target}, learningMode)
}

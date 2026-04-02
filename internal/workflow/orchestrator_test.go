package workflow

import (
	"context"
	"strings"
	"testing"
	"github.com/gentleman-programming/argus/internal/components/engram"
	"github.com/gentleman-programming/argus/internal/components/scanner"
)

// MockProvider simula una IA para pruebas
type MockProvider struct{}

func (m *MockProvider) Chat(ctx context.Context, systemPrompt, userMessage string) (string, error) {
	if strings.Contains(userMessage, "nmap") {
		return `{"action": "TOOL", "category": "Identify", "tool": "nmap", "task": "scan_ports"}`, nil
	}
	return `{"action": "FINISH"}`, nil
}
func (m *MockProvider) ProviderName() string { return "Mock" }
func (m *MockProvider) IsLocal() bool        { return true }

func TestArgusWorkflow(t *testing.T) {
	ctx := context.Background()
	provider := &MockProvider{}
	memory, _ := engram.InjectEngram(ctx, "memory", "")
	
	orch := NewOrchestrator(provider, memory, nil, func(a scanner.ProposedAction) bool { return true })

	t.Run("Autonomous Decision Parsing", func(t *testing.T) {
		// Test parsing and execution flow
		err := orch.RunAutonomousWorkflow(ctx, "localhost")
		if err != nil {
			t.Errorf("Workflow failed: %v", err)
		}
	})
}

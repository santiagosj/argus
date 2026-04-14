package workflow

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gentleman-programming/argus/internal/components/engram"
	"github.com/gentleman-programming/argus/internal/system"
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

	// Mock HITL validator that always approves
	mockValidator := &MockHitlValidator{}

	orch := NewOrchestrator(provider, memory, nil, mockValidator)

	t.Run("Autonomous Decision Parsing", func(t *testing.T) {
		// Test parsing and execution flow
		err := orch.RunAutonomousWorkflow(ctx, "localhost")
		if err != nil {
			t.Errorf("Workflow failed: %v", err)
		}
	})
}

// MockHitlValidator for testing
type MockHitlValidator struct{}

func (m *MockHitlValidator) ProposeAction(ctx context.Context, action system.ProposedAction) (*system.HitlDecision, error) {
	return &system.HitlDecision{
		ApprovedBy:    "TEST",
		ApprovedAt:    time.Now(),
		Reason:        "Auto-approved for testing",
		ActualOutcome: "",
	}, nil
}

func (m *MockHitlValidator) ApproveAction(actionID string, approvedBy string, reason string) error {
	return nil
}

func (m *MockHitlValidator) RejectAction(actionID string, rejectedBy string, reason string) error {
	return nil
}

func (m *MockHitlValidator) GetPendingActions(ctx context.Context) ([]system.ProposedAction, error) {
	return []system.ProposedAction{}, nil
}

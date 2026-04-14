package system

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// InteractiveValidator implements HitlValidator for interactive approval
type InteractiveValidator struct {
	pendingActions map[string]ProposedAction
	mu             sync.RWMutex
	tuiMode        bool // When true, auto-approve actions (user already consented via TUI)
}

// NewInteractiveValidator creates a new interactive validator
func NewInteractiveValidator() *InteractiveValidator {
	return &InteractiveValidator{
		pendingActions: make(map[string]ProposedAction),
		tuiMode:        false,
	}
}

// NewTUIValidator creates a validator for TUI context (auto-approves)
func NewTUIValidator() *InteractiveValidator {
	return &InteractiveValidator{
		pendingActions: make(map[string]ProposedAction),
		tuiMode:        true,
	}
}

// ProposeAction proposes an action for HITL approval
func (v *InteractiveValidator) ProposeAction(ctx context.Context, action ProposedAction) (*HitlDecision, error) {
	// Generate UUID if not provided
	if action.ID == "" {
		action.ID = uuid.New().String()
	}

	// Store pending action
	v.mu.Lock()
	v.pendingActions[action.ID] = action
	v.mu.Unlock()

	// If in TUI mode, auto-approve (user already consented via tool selection)
	if v.tuiMode {
		return &HitlDecision{
			ApprovedBy:    "TUI_USER",
			ApprovedAt:    time.Now(),
			Reason:        "Auto-approved via TUI tool selection",
			ActualOutcome: "",
		}, nil
	}

	// For CLI context, return pending status (would need interactive approval)
	return &HitlDecision{
		ApprovedBy:    "",
		ApprovedAt:    time.Now(),
		Reason:        "Pending user approval in CLI",
		ActualOutcome: "",
	}, fmt.Errorf("action %s requires user approval", action.ID)
}

// ApproveAction approves a pending action
func (v *InteractiveValidator) ApproveAction(actionID string, approvedBy string, reason string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	_, exists := v.pendingActions[actionID]
	if !exists {
		return fmt.Errorf("action %s not found", actionID)
	}

	// Remove from pending
	delete(v.pendingActions, actionID)

	fmt.Printf("✅ Action %s APPROVED by %s: %s\n", actionID, approvedBy, reason)
	return nil
}

// RejectAction rejects a pending action
func (v *InteractiveValidator) RejectAction(actionID string, rejectedBy string, reason string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	_, exists := v.pendingActions[actionID]
	if !exists {
		return fmt.Errorf("action %s not found", actionID)
	}

	// Remove from pending
	delete(v.pendingActions, actionID)

	fmt.Printf("❌ Action %s REJECTED by %s: %s\n", actionID, rejectedBy, reason)
	return nil
}

// GetPendingActions returns all pending actions
func (v *InteractiveValidator) GetPendingActions(ctx context.Context) ([]ProposedAction, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	actions := make([]ProposedAction, 0, len(v.pendingActions))
	for _, action := range v.pendingActions {
		actions = append(actions, action)
	}

	return actions, nil
}

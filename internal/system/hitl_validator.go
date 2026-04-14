package system

import (
	"context"
)

type ProposedAction struct {
	ID           string `json:"id"`           // UUID
	Phase        string `json:"phase"`        // IDENTIFY|PROTECT|etc
	Action       string `json:"action"`       // qué hacer
	Target       string `json:"target"`       // sobre qué
	RiskLevel    string `json:"risk_level"`   // LOW|MEDIUM|HIGH|CRITICAL
	Reason       string `json:"reason"`       // por qué lo propone
	Consequences string `json:"consequences"` // qué puede pasar
}

type HitlValidator interface {
	ProposeAction(ctx context.Context, action ProposedAction) (*HitlDecision, error)
	ApproveAction(actionID string, approvedBy string, reason string) error
	RejectAction(actionID string, rejectedBy string, reason string) error
	GetPendingActions(ctx context.Context) ([]ProposedAction, error)
}

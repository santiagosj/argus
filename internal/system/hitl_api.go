package system

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// HitlAPI provides REST endpoints for HITL operations
type HitlAPI struct {
	hitlValidator HitlValidator
	auditLogger   *AuditLogger
}

// NewHitlAPI creates a new HITL API handler
func NewHitlAPI(hitlValidator HitlValidator, auditLogger *AuditLogger) *HitlAPI {
	return &HitlAPI{
		hitlValidator: hitlValidator,
		auditLogger:   auditLogger,
	}
}

// RegisterRoutes registers the HITL API routes
func (api *HitlAPI) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/hitl/pending", api.GetPendingActions).Methods("GET")
	router.HandleFunc("/api/v1/hitl/approve/{actionID}", api.ApproveAction).Methods("POST")
	router.HandleFunc("/api/v1/hitl/reject/{actionID}", api.RejectAction).Methods("POST")
	router.HandleFunc("/api/v1/hitl/decisions", api.GetDecisions).Methods("GET")
}

// GetPendingActions returns all pending HITL actions
func (api *HitlAPI) GetPendingActions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	actions, err := api.hitlValidator.GetPendingActions(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get pending actions: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"pending_actions": actions,
		"count":           len(actions),
		"timestamp":       time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

// ApproveActionRequest represents an approval request
type ApproveActionRequest struct {
	ApprovedBy string `json:"approved_by"`
	Reason     string `json:"reason"`
}

// ApproveAction approves a pending action
func (api *HitlAPI) ApproveAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	actionID := vars["actionID"]
	if actionID == "" {
		http.Error(w, "actionID is required", http.StatusBadRequest)
		return
	}

	var req ApproveActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.ApprovedBy == "" {
		http.Error(w, "approved_by is required", http.StatusBadRequest)
		return
	}

	err := api.hitlValidator.ApproveAction(actionID, req.ApprovedBy, req.Reason)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to approve action: %v", err), http.StatusInternalServerError)
		return
	}

	// Log the approval
	api.auditLogger.LogWithHitl(
		"HITL-API",
		"HITL_DECISION",
		"APPROVED",
		map[string]string{
			"action_id":   actionID,
			"approved_by": req.ApprovedBy,
			"reason":      req.Reason,
		},
		&HitlDecision{
			ApprovedBy:    req.ApprovedBy,
			ApprovedAt:    time.Now(),
			Reason:        req.Reason,
			ActualOutcome: "",
		},
	)

	response := map[string]interface{}{
		"status":      "approved",
		"action_id":   actionID,
		"approved_by": req.ApprovedBy,
		"timestamp":   time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RejectActionRequest represents a rejection request
type RejectActionRequest struct {
	RejectedBy string `json:"rejected_by"`
	Reason     string `json:"reason"`
}

// RejectAction rejects a pending action
func (api *HitlAPI) RejectAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	actionID := vars["actionID"]
	if actionID == "" {
		http.Error(w, "actionID is required", http.StatusBadRequest)
		return
	}

	var req RejectActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if req.RejectedBy == "" {
		http.Error(w, "rejected_by is required", http.StatusBadRequest)
		return
	}

	err := api.hitlValidator.RejectAction(actionID, req.RejectedBy, req.Reason)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reject action: %v", err), http.StatusInternalServerError)
		return
	}

	// Log the rejection
	api.auditLogger.LogWithHitl(
		"HITL-API",
		"HITL_DECISION",
		"REJECTED",
		map[string]string{
			"action_id":   actionID,
			"rejected_by": req.RejectedBy,
			"reason":      req.Reason,
		},
		&HitlDecision{
			ApprovedBy:    req.RejectedBy,
			ApprovedAt:    time.Now(),
			Reason:        req.Reason,
			ActualOutcome: "",
		},
	)

	response := map[string]interface{}{
		"status":      "rejected",
		"action_id":   actionID,
		"rejected_by": req.RejectedBy,
		"timestamp":   time.Now(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetDecisions returns recent HITL decisions
func (api *HitlAPI) GetDecisions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Implement actual decision retrieval from audit log
	// For now, return empty response
	response := map[string]interface{}{
		"decisions": []interface{}{},
		"count":     0,
		"timestamp": time.Now(),
	}

	json.NewEncoder(w).Encode(response)
}

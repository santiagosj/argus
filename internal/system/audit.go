package system

import (
	"encoding/json"
	"os"
	"time"
)

type AuditEntry struct {
	Timestamp time.Time     `json:"timestamp"`
	Source    string        `json:"source"`
	Type      string        `json:"type"` // THINK, ACTION, FINDING, ERROR, PROPOSAL, HITL_REQUEST, HITL_DECISION
	Content   interface{}   `json:"content"`
	Status    string        `json:"status"`              // PENDING, APPROVED, REJECTED, EXECUTED
	HitlData  *HitlDecision `json:"hitl_data,omitempty"` // null si no requiere aprobación
}

type AuditLogger struct {
	file *os.File
}

type HitlDecision struct {
	ApprovedBy    string    // Usuario que aprobó
	ApprovedAt    time.Time // Timestamp
	Reason        string    // Por qué aprobó/rechazó
	ActualOutcome string    // Qué pasó realmente vs predicho
}

func NewAuditLogger(path string) (*AuditLogger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &AuditLogger{file: f}, nil
}

func (l *AuditLogger) Log(source, entryType string, content interface{}) {
	entry := AuditEntry{
		Timestamp: time.Now(),
		Source:    source,
		Type:      entryType,
		Content:   content,
		Status:    "EXECUTED", // Default for backward compatibility
		HitlData:  nil,
	}

	jsonData, _ := json.Marshal(entry)
	l.file.Write(jsonData)
	l.file.WriteString("\n")
}

// LogWithHitl logs an entry with HITL validation information
func (l *AuditLogger) LogWithHitl(source, entryType, status string, content interface{}, hitlData *HitlDecision) {
	entry := AuditEntry{
		Timestamp: time.Now(),
		Source:    source,
		Type:      entryType,
		Content:   content,
		Status:    status,
		HitlData:  hitlData,
	}

	jsonData, _ := json.Marshal(entry)
	l.file.Write(jsonData)
	l.file.WriteString("\n")
}

func (l *AuditLogger) Close() {
	l.file.Close()
}

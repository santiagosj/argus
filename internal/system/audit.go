package system

import (
	"encoding/json"
	"os"
	"time"
)

type AuditEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	Source    string      `json:"source"`
	Type      string      `json:"type"` // THINK, ACTION, FINDING, ERROR
	Content   interface{} `json:"content"`
}

type AuditLogger struct {
	file *os.File
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
	}

	jsonData, _ := json.Marshal(entry)
	l.file.Write(jsonData)
	l.file.WriteString("\n")
}

func (l *AuditLogger) Close() {
	l.file.Close()
}

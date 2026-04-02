package engram

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// MemoryStore representa la abstracción de Engram para ciberseguridad.
type MemoryStore interface {
	SaveFinding(findingID string, data map[string]interface{}) error
	RetrieveContext(lastHours int) ([]string, error)
	Cleanup() error
}

// MemorySession captura la "historia" de un escaneo o auditoría.
type MemorySession struct {
	ID        string
	Project   string
	StartTime time.Time
}

// InjectEngram configura el agente de Argus con persistencia cognitiva.
func InjectEngram(ctx context.Context, storageType, path string) (MemoryStore, error) {
	if storageType == "sqlite" {
		return NewSQLiteMemory(path)
	}
	return &SimpleMemory{findings: make(map[string]map[string]interface{})}, nil
}

type SQLiteMemory struct {
	db *sql.DB
}

func NewSQLiteMemory(path string) (*SQLiteMemory, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS findings (
		id TEXT PRIMARY KEY,
		data TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}

	return &SQLiteMemory{db: db}, nil
}

func (s *SQLiteMemory) SaveFinding(findingID string, data map[string]interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT OR REPLACE INTO findings (id, data) VALUES (?, ?)", findingID, string(jsonData))
	return err
}

func (s *SQLiteMemory) RetrieveContext(lastHours int) ([]string, error) {
	rows, err := s.db.Query("SELECT id, data FROM findings WHERE timestamp >= datetime('now', ?)", fmt.Sprintf("-%d hours", lastHours))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var context []string
	for rows.Next() {
		var id, data string
		if err := rows.Scan(&id, &data); err != nil {
			continue
		}
		context = append(context, fmt.Sprintf("%s: %s", id, data))
	}
	return context, nil
}

func (s *SQLiteMemory) Cleanup() error {
	_, err := s.db.Exec("DELETE FROM findings")
	return err
}

type SimpleMemory struct {
	findings map[string]map[string]interface{}
}

func (s *SimpleMemory) SaveFinding(findingID string, data map[string]interface{}) error {
	s.findings[findingID] = data
	return nil
}

func (s *SimpleMemory) RetrieveContext(lastHours int) ([]string, error) {
	var context []string
	for k, v := range s.findings {
		context = append(context, fmt.Sprintf("%s: %v", k, v))
	}
	return context, nil
}

func (s *SimpleMemory) Cleanup() error {
	s.findings = make(map[string]map[string]interface{})
	return nil
}

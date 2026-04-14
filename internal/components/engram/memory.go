package engram

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
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

type Entity struct {
	Type  string
	Value string
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
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		ttl_hours INTEGER DEFAULT 24
	)`)
	if err != nil {
		return nil, err
	}

	// Create entities table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS entities (
		id INTEGER PRIMARY KEY,
		finding_id TEXT,
		entity_type TEXT,
		entity_value TEXT,
		FOREIGN KEY (finding_id) REFERENCES findings(id)
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
	if err != nil {
		return err
	}

	// Extract and save entities
	entities := s.extractEntities(string(jsonData))
	for _, entity := range entities {
		_, err = s.db.Exec("INSERT INTO entities (finding_id, entity_type, entity_value) VALUES (?, ?, ?)", findingID, entity.Type, entity.Value)
		if err != nil {
			// Log but don't fail
			fmt.Printf("Error saving entity: %v\n", err)
		}
	}

	return nil
}

func (s *SQLiteMemory) RetrieveContext(lastHours int) ([]string, error) {
	rows, err := s.db.Query("SELECT id, data FROM findings WHERE timestamp >= datetime('now', ?) AND (ttl_hours IS NULL OR timestamp >= datetime('now', '-' || ttl_hours || ' hours'))", fmt.Sprintf("-%d hours", lastHours))
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
	// Delete expired findings
	_, err := s.db.Exec("DELETE FROM findings WHERE ttl_hours IS NOT NULL AND timestamp < datetime('now', '-' || ttl_hours || ' hours')")
	if err != nil {
		return err
	}
	// Delete orphaned entities
	_, err = s.db.Exec("DELETE FROM entities WHERE finding_id NOT IN (SELECT id FROM findings)")
	return err
}

func (s *SQLiteMemory) extractEntities(text string) []Entity {
	var entities []Entity

	// IP addresses
	ipRegex := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	for _, match := range ipRegex.FindAllString(text, -1) {
		entities = append(entities, Entity{Type: "IP", Value: match})
	}

	// URLs
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	for _, match := range urlRegex.FindAllString(text, -1) {
		entities = append(entities, Entity{Type: "URL", Value: match})
	}

	// CVEs
	cveRegex := regexp.MustCompile(`CVE-\d{4}-\d{4,}`)
	for _, match := range cveRegex.FindAllString(text, -1) {
		entities = append(entities, Entity{Type: "CVE", Value: match})
	}

	// CWEs
	cweRegex := regexp.MustCompile(`CWE-\d+`)
	for _, match := range cweRegex.FindAllString(text, -1) {
		entities = append(entities, Entity{Type: "CWE", Value: match})
	}

	return entities
}

func (s *SQLiteMemory) FindSimilarFindings(query string) ([]string, error) {
	entities := s.extractEntities(query)
	if len(entities) == 0 {
		return []string{}, nil
	}

	var conditions []string
	var args []interface{}
	for _, entity := range entities {
		conditions = append(conditions, "(entity_type = ? AND entity_value = ?)")
		args = append(args, entity.Type, entity.Value)
	}

	querySQL := fmt.Sprintf("SELECT DISTINCT f.id, f.data FROM findings f JOIN entities e ON f.id = e.finding_id WHERE %s", strings.Join(conditions, " OR "))
	rows, err := s.db.Query(querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var id, data string
		if err := rows.Scan(&id, &data); err != nil {
			continue
		}
		results = append(results, fmt.Sprintf("%s: %s", id, data))
	}
	return results, nil
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

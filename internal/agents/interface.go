package agents

import (
	"context"
)

// NISTCategory representa las 5 funciones del NIST CSF.
type NISTCategory string

const (
	Identify NISTCategory = "Identify"
	Protect  NISTCategory = "Protect"
	Detect   NISTCategory = "Detect"
	Respond  NISTCategory = "Respond"
	Recover  NISTCategory = "Recover"
	Evolve   NISTCategory = "Evolve"
	General  NISTCategory = "General"
)

// AgentProvider es la abstracción para conectarse con IA (Ollama, Claude, etc).
type AgentProvider interface {
	// Identity
	ProviderName() string
	
	// Thinking - Prompting core
	Chat(ctx context.Context, systemPrompt string, userMessage string) (string, error)
	
	// Local/Cloud
	IsLocal() bool
}

// CyberAgent es el orquestador cognitivo que realiza tareas de seguridad.
type CyberAgent interface {
	// NIST mapping
	Category() NISTCategory
	
	// Skills
	LoadSkills(skillNames []string) error
	
	// Context management (Engram integration)
	WithMemory(memoryID string) CyberAgent
	
	// Execution
	ExecuteTask(ctx context.Context, task string) (Report, error)
}

// Report representa el resultado de una tarea (hallazgos, severidad, remediación).
type Report struct {
	Summary     string
	Findings    []Finding
	Severity    string // Low, Medium, High, Critical
	Remediation string
}

type Finding struct {
	Tool     string
	Evidence string
	CWE      string
}

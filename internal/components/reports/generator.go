package reports

import (
	"context"
	"fmt"

	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/components/engram"
)

// ReportGenerator consolida hallazgos en un reporte final.
type ReportGenerator struct {
	provider agents.AgentProvider
	memory   engram.MemoryStore
}

func NewReportGenerator(p agents.AgentProvider, m engram.MemoryStore) *ReportGenerator {
	return &ReportGenerator{provider: p, memory: m}
}

// GenerateFinalReport crea un reporte Markdown alineado con NIST CSF.
func (g *ReportGenerator) GenerateFinalReport(ctx context.Context, target string) (string, error) {
	// 1. Recuperar todos los hallazgos de Engram
	findings, err := g.memory.RetrieveContext(10) // Traer los últimos 10 hallazgos
	if err != nil {
		return "", err
	}

	systemPrompt := `You are a Senior Cyber-Security Reporter. 
Your goal is to consolidate technical findings into a professional Markdown report.
The report MUST be aligned with the NIST CSF (Identify, Protect, Detect, Respond, Recover).`

	userPrompt := fmt.Sprintf(`Target: %s
Findings from Sub-Agents:
%v

Generate a report with:
1. Executive Summary
2. Technical Findings (grouped by NIST Category)
3. Risk Assessment
4. Recommended Actions`, target, findings)

	report, err := g.provider.Chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	return report, nil
}

package skills

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/argus/internal/agents"
)

// SkillGenerator es el agente de evolución que crea nuevas habilidades.
type SkillGenerator struct {
	provider agents.AgentProvider
}

func NewSkillGenerator(p agents.AgentProvider) *SkillGenerator {
	return &SkillGenerator{provider: p}
}

// GenerateSkillFromDescription crea un archivo SKILL.md basado en el conocimiento del modelo o documentación provista.
func (g *SkillGenerator) GenerateSkillFromDescription(ctx context.Context, category agents.NISTCategory, toolName string, description string) (string, error) {
	systemPrompt := `You are a Argus Evolution Agent. Your goal is to create a professional SKILL.md for a security tool.
The skill must be compatible with the NIST CSF framework.
Output ONLY the Markdown content.`

	userPrompt := fmt.Sprintf(`Create a SKILL.md for the tool "%s" in the NIST category "%s".
Description/Context: %s

The Markdown should include:
1. Tool Description
2. Common Commands (one per line)
3. How to interpret output
4. Security best practices`, toolName, category, description)

	content, err := g.provider.Chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		return "", err
	}

	// Guardar el archivo
	dirPath := filepath.Join("skills", strings.ToLower(string(category)))
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", err
	}

	filePath := filepath.Join(dirPath, strings.ToLower(toolName)+".md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", err
	}

	return filePath, nil
}

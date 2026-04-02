package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/components/engram"
	"github.com/gentleman-programming/argus/internal/components/reports"
	"github.com/gentleman-programming/argus/internal/components/scanner"
	"github.com/gentleman-programming/argus/internal/components/skills"
	"github.com/gentleman-programming/argus/internal/system"
)

// Orchestrator gestiona el equipo de sub-agentes.
type Orchestrator struct {
	provider  agents.AgentProvider
	memory    engram.MemoryStore
	evolution *skills.SkillGenerator
	scanner   *scanner.CommandExecutor
	dashboard *system.WebDashboard
	validator func(scanner.ProposedAction) bool
	reporter  *reports.ReportGenerator
	audit     *system.AuditLogger
}

func NewOrchestrator(p agents.AgentProvider, m engram.MemoryStore, d *system.WebDashboard, v func(scanner.ProposedAction) bool) *Orchestrator {
	audit, _ := system.NewAuditLogger("argus_audit.jsonl")
	return &Orchestrator{
		provider:  p,
		memory:    m,
		evolution: skills.NewSkillGenerator(p),
		scanner:   &scanner.CommandExecutor{},
		dashboard: d,
		validator: v,
		reporter:  reports.NewReportGenerator(p, m),
		audit:     audit,
	}
}

func (o *Orchestrator) SpawnSubOrchestrator(ctx context.Context, target string, task string) error {
	o.logToDashboard("HIERARCHY", fmt.Sprintf("Spawning Sub-Orchestrator for task: %s on target: %s", task, target))
	o.audit.Log("Argus-Lead", "SPAWN", map[string]string{"target": target, "task": task})

	subOrch := NewOrchestrator(o.provider, o.memory, o.dashboard, o.validator)
	return subOrch.RunAutonomousWorkflow(ctx, target)
}

type AgentDecision struct {
	Action    string `json:"action"` // TOOL, SUB_ORCHESTRATOR, FINISH
	Tool      string `json:"tool"`
	Cat       string `json:"category"`
	Task      string `json:"task"`
	SubTarget string `json:"sub_target"`
}

// RunAutonomousWorkflow es el cerebro dinámico que decide qué hacer sin pasos hardcodeados.
func (o *Orchestrator) RunAutonomousWorkflow(ctx context.Context, target string) error {
	o.logToDashboard("STATUS", fmt.Sprintf("Argus Engine started for: %s", target))
	o.audit.Log("Argus-Lead", "START", target)

	if strings.Contains(target, "evolve") || strings.Contains(target, "self") {
		return o.RunEvolutionWorkflow(ctx, target)
	}

	for i := 0; i < 5; i++ {
		o.logToDashboard("BRAIN", "Thinking about the next logical step...")

		findings, _ := o.memory.RetrieveContext(1)
		systemPrompt := `You are the Lead Pentester Agent for Argus. 
Decide the next step based on findings. 
If the task is complex, you can spawn a Sub-Orchestrator for a specific target or subdomain.
Output format: JSON ONLY
{
  "action": "TOOL" | "SUB_ORCHESTRATOR" | "FINISH",
  "tool": "name",
  "category": "NIST_CAT",
  "task": "description",
  "sub_target": "if spawning sub-orchestrator"
}`

		userPrompt := fmt.Sprintf("Current Context: %s\nCurrent Findings: %v", target, findings)
		response, err := o.provider.Chat(ctx, systemPrompt, userPrompt)
		if err != nil {
			o.audit.Log("Argus-Lead", "ERROR", err.Error())
			break
		}

		o.audit.Log("Argus-Lead", "THINK", response)

		var decision AgentDecision
		err = json.Unmarshal([]byte(extractJSON(response)), &decision)
		if err != nil {
			o.logToDashboard("ERROR", "Failed to parse brain decision")
			continue
		}

		if decision.Action == "FINISH" {
			o.audit.Log("Argus-Lead", "FINISH", "Goal reached")
			break
		}

		if decision.Action == "SUB_ORCHESTRATOR" {
			subTarget := decision.SubTarget
			if subTarget == "" {
				subTarget = target
			}
			o.SpawnSubOrchestrator(ctx, subTarget, decision.Task)
			continue
		}

		tool := decision.Tool
		cat := decision.Cat
		task := decision.Task

		if tool == "" {
			continue
		}

		if !system.CheckTool(tool) {
			o.logToDashboard("SYSTEM", fmt.Sprintf("Tool %s missing. Proposing installation...", tool))
			o.handleMissingTool(ctx, tool)
			if !system.CheckTool(tool) {
				continue
			}
		}

		skillPath := filepath.Join("skills", cat, tool+".md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			o.logToDashboard("EVOLUTION", fmt.Sprintf("Learning tool %s...", tool))
			o.LearnNewTool(ctx, agents.NISTCategory(cat), tool, "Auto-discovered tool")
		}

		o.logToDashboard("STATUS", fmt.Sprintf("Executing: %s using %s", task, tool))
		_, err = o.runSubAgentActionable(ctx, agents.NISTCategory(cat), tool, task)
		if err != nil {
			o.logToDashboard("ERROR", fmt.Sprintf("Execution failed: %v. Analyzing why...", err))
			o.audit.Log("Argus-Lead", "EXEC_ERROR", err.Error())
			// Auto-correction loop could go here
		}
	}

	o.logToDashboard("STATUS", "Argus Workflow finished.")
	report, err := o.reporter.GenerateFinalReport(ctx, target)
	if err == nil {
		o.logToDashboard("REPORT", report)
		fmt.Printf("\n--- FINAL REPORT ---\n%s\n", report)
	}

	return nil
}

func (o *Orchestrator) runSubAgentActionable(ctx context.Context, cat agents.NISTCategory, skillName string, input string) (*agents.Report, error) {
	o.logToDashboard("AGENT", fmt.Sprintf("Thinking about: %s", input))

	skill, err := LoadSkill(string(cat), skillName)
	if err != nil {
		return nil, err
	}
	systemPrompt := fmt.Sprintf("You are a specialized %s security agent. \nSkill Context: %s\nOutput ONLY the command or 'NO_COMMAND'.", cat, skill.Content)

	contextInfo, _ := o.memory.RetrieveContext(1)
	fullPrompt := fmt.Sprintf("Findings: %v\nTask: %s", contextInfo, input)

	command, err := o.provider.Chat(ctx, systemPrompt, fullPrompt)
	if err != nil || strings.Contains(command, "NO_COMMAND") {
		return o.runSubAgent(ctx, cat, skillName, input)
	}

	action := scanner.ProposedAction{
		Command:     strings.TrimSpace(command),
		Description: fmt.Sprintf("Agent %s using %s", cat, skillName),
	}

	output, err := o.scanner.ValidateAndRun(ctx, action, o.validator)
	if err != nil {
		// En Argus, devolvemos el error enriquecido para que el orquestador lo analice
		return nil, err
	}

	o.memory.SaveFinding(skillName, map[string]interface{}{"output": output, "command": action.Command})
	o.logToDashboard(skillName, output)
	o.audit.Log("Argus-Agent", "ACTION", map[string]string{"tool": skillName, "output": output})

	return &agents.Report{Summary: output}, nil
}

func (o *Orchestrator) RunEvolutionWorkflow(ctx context.Context, input string) error {
	o.logToDashboard("STATUS", "Evolution Cycle (SDD) Started")
	o.audit.Log("Argus-Evolve", "START", input)
	changeName := fmt.Sprintf("evolve-%d", os.Getpid())

	// Reducido a solo las fases críticas para mayor velocidad
	phases := []string{"sdd-explore", "sdd-propose", "sdd-apply"}
	for _, phase := range phases {
		if phase == "sdd-apply" {
			o.logToDashboard("EVOLVE", "Applying changes...")
			applyReport, err := o.runSubAgentApply(ctx, changeName)
			if err != nil {
				return err
			}
			o.audit.Log("Argus-Evolve", "APPLY", applyReport.Summary)
			continue
		}

		o.logToDashboard("EVOLVE", fmt.Sprintf("Executing %s...", phase))
		report, err := o.runSubAgent(ctx, agents.Evolve, phase, input)
		if err != nil {
			return err
		}
		o.audit.Log("Argus-Evolve", "PHASE", map[string]string{"phase": phase, "report": report.Summary})
	}

	o.logToDashboard("STATUS", "Evolution Cycle Finished. Argus has evolved.")
	return nil
}

func (o *Orchestrator) runSubAgentApply(ctx context.Context, changeName string) (*agents.Report, error) {
	o.logToDashboard("EVOLVE", "Executing SDD-Apply sub-agent...")
	skill, err := LoadSkill("Evolve", "sdd-apply")
	if err != nil {
		return nil, err
	}
	common, err := LoadSkill("_shared", "sdd-phase-common")
	if err != nil {
		return nil, err
	}

	systemPrompt := fmt.Sprintf("You are the SDD-Apply agent. Protocol: %s\nSkill: %s\nOutput JSON ProposedAction.", common.Content, skill.Content)
	contextInfo, _ := o.memory.RetrieveContext(5)
	fullPrompt := fmt.Sprintf("Change: %s\nContext: %v", changeName, contextInfo)

	response, err := o.provider.Chat(ctx, systemPrompt, fullPrompt)
	if err != nil {
		return nil, err
	}

	var action scanner.ProposedAction
	err = json.Unmarshal([]byte(extractJSON(response)), &action)
	if err != nil {
		action = scanner.ProposedAction{Type: scanner.ActionCommand, Command: strings.TrimSpace(response), Description: "Evolution command"}
	}

	output, err := o.scanner.ValidateAndRun(ctx, action, o.validator)
	if err != nil {
		return nil, err
	}

	return &agents.Report{Summary: output}, nil
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || end < start {
		return s
	}
	return s[start : end+1]
}

func (o *Orchestrator) handleMissingTool(ctx context.Context, tool string) {
	action := scanner.ProposedAction{
		Command:     fmt.Sprintf("sudo apt install -y %s", tool),
		Description: fmt.Sprintf("Install missing security tool: %s", tool),
	}
	o.scanner.ValidateAndRun(ctx, action, o.validator)
}

func (o *Orchestrator) logToDashboard(source string, content interface{}) {
	if o.dashboard != nil {
		o.dashboard.Broadcast(system.DashboardUpdate{
			Type:    "finding",
			Source:  source,
			Content: content,
		})
	}
}

func (o *Orchestrator) LearnNewTool(ctx context.Context, cat agents.NISTCategory, tool string, desc string) error {
	path, err := o.evolution.GenerateSkillFromDescription(ctx, cat, tool, desc)
	if err != nil {
		return err
	}
	o.audit.Log("Argus-Evolve", "LEARN", map[string]string{"tool": tool, "path": path})
	return nil
}

func (o *Orchestrator) runSubAgent(ctx context.Context, cat agents.NISTCategory, skillName string, input string) (*agents.Report, error) {
	skill, err := LoadSkill(string(cat), skillName)
	if err != nil {
		return nil, err
	}
	content := skill.Content
	if cat == agents.Evolve && strings.HasPrefix(skillName, "sdd-") {
		common, err := LoadSkill("_shared", "sdd-phase-common")
		if err == nil {
			content = common.Content + "\n\n" + content
		}
	}

	systemPrompt := fmt.Sprintf("You are a specialized %s security agent. \nSkill Context: %s", cat, content)
	contextInfo, _ := o.memory.RetrieveContext(1)
	fullPrompt := fmt.Sprintf("Findings: %v\nTask: %s", contextInfo, input)

	response, err := o.provider.Chat(ctx, systemPrompt, fullPrompt)
	if err != nil {
		return nil, err
	}

	o.memory.SaveFinding(skillName, map[string]interface{}{"output": response})
	return &agents.Report{Summary: response}, nil
}

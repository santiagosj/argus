package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/gentleman-programming/argus/internal/agents"
	"github.com/gentleman-programming/argus/internal/components/engram"
	"github.com/gentleman-programming/argus/internal/components/reports"
	"github.com/gentleman-programming/argus/internal/components/scanner"
	"github.com/gentleman-programming/argus/internal/components/skills"
	"github.com/gentleman-programming/argus/internal/system"
)

// Color styles
var colorBold = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")).Bold(true)
var colorInfo = lipgloss.NewStyle().Foreground(lipgloss.Color("#00d4ff"))
var colorWarn = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffaa00"))
var colorError = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5555"))
var colorSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("#55ff55"))

var spinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// Orchestrator gestiona el equipo de sub-agentes.
type Orchestrator struct {
	provider  agents.AgentProvider
	memory    engram.MemoryStore
	evolution *skills.SkillGenerator
	scanner   *scanner.CommandExecutor
	dashboard *system.WebDashboard
	validator system.HitlValidator // HITL validator for human approval
	reporter  *reports.ReportGenerator
	audit     *system.AuditLogger
	mu        sync.Mutex
}

func NewOrchestrator(p agents.AgentProvider, m engram.MemoryStore, d *system.WebDashboard, v system.HitlValidator) *Orchestrator {
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

type SDDPhaseResult struct {
	Phase            string
	Status           string
	ExecutiveSummary string
	DetailedReport   string
	Artifacts        []string
	NextRecommended  string
	Risks            string
	SkillResolution  string
	Raw              string
}

// RunAutonomousWorkflow es el cerebro dinámico que decide qué hacer sin pasos hardcodeados.
func (o *Orchestrator) RunAutonomousWorkflow(ctx context.Context, target string) error {
	o.logToDashboard("STATUS", fmt.Sprintf("Argus Engine started for: %s", target))
	o.audit.Log("Argus-Lead", "START", target)
	fmt.Printf("Starting Argus workflow on target: %s\n", target)

	if strings.Contains(target, "evolve") || strings.Contains(target, "self") {
		return o.RunEvolutionWorkflow(ctx, target)
	}

	for i := 0; i < 5; i++ {
		fmt.Printf("Iteration %d/5: Analyzing next step...\n", i+1)
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
			fmt.Printf("Error in decision making: %v\n", err)
			break
		}

		o.audit.Log("Argus-Lead", "THINK", response)

		var decision AgentDecision
		err = json.Unmarshal([]byte(extractJSON(response)), &decision)
		if err != nil {
			o.logToDashboard("ERROR", "Failed to parse brain decision")
			fmt.Printf("Failed to parse decision, retrying...\n")
			continue
		}

		if decision.Action == "FINISH" {
			o.audit.Log("Argus-Lead", "FINISH", "Goal reached")
			fmt.Printf("Workflow completed successfully.\n")
			break
		}

		if decision.Action == "SUB_ORCHESTRATOR" {
			subTarget := decision.SubTarget
			if subTarget == "" {
				subTarget = target
			}
			fmt.Printf("Spawning sub-orchestrator for task: %s\n", decision.Task)
			o.SpawnSubOrchestrator(ctx, subTarget, decision.Task)
			continue
		}

		tool := decision.Tool
		cat := decision.Cat
		task := decision.Task

		if tool == "" {
			fmt.Printf("No tool specified, skipping...\n")
			continue
		}

		if !system.CheckTool(tool) {
			o.logToDashboard("SYSTEM", fmt.Sprintf("Tool %s missing. Proposing installation...", tool))
			fmt.Printf("Tool %s not found, attempting installation...\n", tool)
			o.handleMissingTool(ctx, tool)
			if !system.CheckTool(tool) {
				fmt.Printf("Installation failed, skipping tool.\n")
				continue
			}
		}

		skillPath := filepath.Join("skills", cat, tool+".md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			o.logToDashboard("EVOLUTION", fmt.Sprintf("Learning tool %s...", tool))
			fmt.Printf("Learning new skill for tool: %s\n", tool)
			o.LearnNewTool(ctx, agents.NISTCategory(cat), tool, "Auto-discovered tool")
		}

		o.logToDashboard("STATUS", fmt.Sprintf("Executing: %s using %s", task, tool))
		fmt.Printf("Executing tool: %s for task: %s\n", tool, task)
		_, err = o.runSubAgentActionable(ctx, agents.NISTCategory(cat), tool, task)
		if err != nil {
			o.logToDashboard("ERROR", fmt.Sprintf("Execution failed: %v. Analyzing why...", err))
			o.audit.Log("Argus-Lead", "EXEC_ERROR", err.Error())
			fmt.Printf("Execution failed: %v\n", err)
			// Error recovery: try alternative tool
			if altTool := o.findAlternativeTool(ctx, agents.NISTCategory(cat), tool, task, err); altTool != "" {
				fmt.Printf("Trying alternative tool: %s\n", altTool)
				_, altErr := o.runSubAgentActionable(ctx, agents.NISTCategory(cat), altTool, task)
				if altErr == nil {
					fmt.Printf("Alternative tool succeeded.\n")
				} else {
					fmt.Printf("Alternative tool also failed: %v\n", altErr)
				}
			}
		} else {
			fmt.Printf("Execution completed successfully.\n")
		}
	}

	o.logToDashboard("STATUS", "Argus Workflow finished.")
	fmt.Printf("Workflow finished. Generating report...\n")
	report, err := o.reporter.GenerateFinalReport(ctx, target)
	if err == nil {
		o.logToDashboard("REPORT", report)
		fmt.Printf("\n--- FINAL REPORT ---\n%s\n", report)
	}

	return nil
}

func (o *Orchestrator) RunToolOnTarget(ctx context.Context, target, tool string) error {
	category := mapToolToCategory(tool)
	skillPath := filepath.Join("skills", string(category), tool+".md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		o.logToDashboard("EVOLUTION", fmt.Sprintf("Generating skill for %s...", tool))
		fmt.Printf("Generating new skill for tool: %s\n", tool)
		if err := o.LearnNewTool(ctx, category, tool, fmt.Sprintf("Auto-generated skill for %s execution on %s", tool, target)); err != nil {
			o.logToDashboard("ERROR", fmt.Sprintf("Failed to generate skill for %s: %v", tool, err))
			return err
		}
	}

	if !system.CheckTool(tool) {
		o.logToDashboard("SYSTEM", fmt.Sprintf("Tool %s missing. Proposing installation...", tool))
		o.handleMissingTool(ctx, tool)
		if !system.CheckTool(tool) {
			return fmt.Errorf("tool %s not installed", tool)
		}
	}

	task := fmt.Sprintf("Execute %s against %s and return structured scan output", tool, target)
	report, err := o.runSubAgentActionable(ctx, category, tool, task)
	if err != nil {
		return err
	}

	if strings.EqualFold(tool, "nmap") {
		parsed := o.parseNmapOutput(report.Summary)
		if parsed != nil {
			o.logToDashboard("NMAP", parsed)
			return nil
		}
	}

	o.logToDashboard(strings.ToUpper(tool), map[string]interface{}{
		"target": target,
		"output": report.Summary,
	})
	return nil
}

func (o *Orchestrator) RunToolOnTargetWithUpdates(ctx context.Context, target, tool string, outputChan chan<- string) error {
	// Send initial status
	outputChan <- fmt.Sprintf("🔧 Starting %s scan on %s...", tool, target)

	category := mapToolToCategory(tool)
	skillPath := filepath.Join("skills", string(category), tool+".md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		outputChan <- fmt.Sprintf("📚 Generating skill for %s...", tool)
		o.logToDashboard("EVOLUTION", fmt.Sprintf("Generating skill for %s...", tool))
		fmt.Printf("Generating new skill for tool: %s\n", tool)
		if err := o.LearnNewTool(ctx, category, tool, fmt.Sprintf("Auto-generated skill for %s execution on %s", tool, target)); err != nil {
			outputChan <- fmt.Sprintf("❌ Failed to generate skill for %s: %v", tool, err)
			o.logToDashboard("ERROR", fmt.Sprintf("Failed to generate skill for %s: %v", tool, err))
			return err
		}
		outputChan <- fmt.Sprintf("✅ Skill generated for %s", tool)
	}

	if !system.CheckTool(tool) {
		outputChan <- fmt.Sprintf("📦 Installing missing tool: %s...", tool)
		o.logToDashboard("SYSTEM", fmt.Sprintf("Tool %s missing. Proposing installation...", tool))
		o.handleMissingTool(ctx, tool)
		if !system.CheckTool(tool) {
			outputChan <- fmt.Sprintf("❌ Tool %s not installed", tool)
			return fmt.Errorf("tool %s not installed", tool)
		}
		outputChan <- fmt.Sprintf("✅ Tool %s installed", tool)
	}

	outputChan <- fmt.Sprintf("⚡ Executing %s against %s...", tool, target)
	task := fmt.Sprintf("Execute %s against %s and return structured scan output", tool, target)
	report, err := o.runSubAgentActionable(ctx, category, tool, task)
	if err != nil {
		outputChan <- fmt.Sprintf("❌ Execution failed: %v", err)
		return err
	}

	outputChan <- "📊 Processing results..."
	if strings.EqualFold(tool, "nmap") {
		parsed := o.parseNmapOutput(report.Summary)
		if parsed != nil {
			outputChan <- fmt.Sprintf("✅ Nmap scan completed. Found %d ports/services", len(parsed))
			o.logToDashboard("NMAP", parsed)
		}
	}

	outputChan <- fmt.Sprintf("✅ %s execution completed successfully", tool)
	o.logToDashboard(strings.ToUpper(tool), map[string]interface{}{
		"target": target,
		"output": report.Summary,
	})
	return nil
}

func mapToolToCategory(tool string) agents.NISTCategory {
	switch strings.ToLower(tool) {
	case "nmap", "feroxbuster", "gobuster", "dirb":
		return agents.Identify
	case "nuclei":
		return agents.Detect
	default:
		return agents.General
	}
}

func (o *Orchestrator) parseNmapOutput(output string) map[string]interface{} {
	result := map[string]interface{}{"raw": output}
	var ports []map[string]string
	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Nmap scan report for") {
			result["host"] = strings.TrimSpace(strings.TrimPrefix(trimmed, "Nmap scan report for"))
			continue
		}
		if strings.HasPrefix(trimmed, "OS details:") {
			result["os"] = strings.TrimSpace(strings.TrimPrefix(trimmed, "OS details:"))
			continue
		}
		if strings.HasPrefix(trimmed, "Service Info:") {
			result["service_info"] = strings.TrimSpace(strings.TrimPrefix(trimmed, "Service Info:"))
			continue
		}
		if strings.HasPrefix(trimmed, "PORT") {
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) >= 3 && strings.Contains(fields[0], "/") {
			ports = append(ports, map[string]string{
				"port":    fields[0],
				"state":   fields[1],
				"service": fields[2],
				"info":    strings.Join(fields[3:], " "),
			})
		}
	}
	if len(ports) > 0 {
		result["ports"] = ports
	}
	return result
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
		Type:        scanner.ActionCommand,
		Command:     strings.TrimSpace(command),
		Description: fmt.Sprintf("Agent %s using %s", cat, skillName),
	}

	output, err := o.scanner.ValidateAndRunWithHitl(ctx, action, o.validator, string(cat))
	if err != nil {
		// En Argus, devolvemos el error enriquecido para que el orquestador lo analice
		return nil, err
	}

	o.memory.SaveFinding(skillName, map[string]interface{}{"output": output, "command": action.Command})
	o.logToDashboard(skillName, output)
	o.audit.LogWithHitl("Argus-Agent", "ACTION", "EXECUTED", map[string]string{"tool": skillName, "output": output}, nil)

	return &agents.Report{Summary: output}, nil
}

func (o *Orchestrator) RunEvolutionWorkflow(ctx context.Context, input string) error {
	o.logToDashboard("STATUS", "Evolution Cycle (SDD) Started")
	o.audit.Log("Argus-Evolve", "START", input)
	changeName := fmt.Sprintf("evolve-%d", os.Getpid())

	defaultPhases := []string{"sdd-init", "sdd-explore", "sdd-propose", "sdd-spec", "sdd-design", "sdd-tasks", "sdd-apply", "sdd-verify", "sdd-archive"}
	currentPhase := "sdd-init"
	executed := make(map[string]bool)
	project := strings.TrimSpace(input)
	if project == "" {
		project = "default"
	}

	for {
		if currentPhase == "" || strings.EqualFold(currentPhase, "none") {
			break
		}

		if executed[currentPhase] {
			return fmt.Errorf("SDD cycle loop detected at %s", currentPhase)
		}
		executed[currentPhase] = true

		result, err := o.runSDDPhase(ctx, changeName, project, currentPhase, input)
		if err != nil {
			return err
		}

		if strings.EqualFold(result.Status, "blocked") || strings.EqualFold(result.Status, "partial") {
			o.logToDashboard("EVOLVE", fmt.Sprintf("SDD phase %s returned status %s", currentPhase, result.Status))
			return fmt.Errorf("SDD blocked or partial at %s: %s", currentPhase, result.ExecutiveSummary)
		}

		o.audit.Log("Argus-Evolve", "PHASE", map[string]string{"phase": currentPhase, "summary": result.ExecutiveSummary})

		next := normalizeSDDPhase(result.NextRecommended)
		if next == "" || strings.EqualFold(next, "none") {
			next = nextSDDPhase(currentPhase, defaultPhases)
		}
		if next == "" || strings.EqualFold(next, "none") {
			break
		}

		currentPhase = next
	}

	o.logToDashboard("STATUS", "Evolution Cycle Finished. Argus has evolved.")
	return nil
}

func (o *Orchestrator) runSubAgentApply(ctx context.Context, changeName string) (*SDDPhaseResult, error) {
	o.logToDashboard("EVOLVE", "Executing SDD-Apply sub-agent...")
	skill, err := LoadSkill("Evolve", "sdd-apply")
	if err != nil {
		return nil, err
	}
	common, err := LoadSkill("_shared", "sdd-phase-common")
	if err != nil {
		return nil, err
	}

	systemPrompt := fmt.Sprintf("You are the SDD-Apply agent. Protocol: %s\nSkill: %s\nOutput JSON ProposedAction or a structured SDD envelope.", common.Content, skill.Content)
	contextInfo, _ := o.memory.RetrieveContext(48)
	fullPrompt := fmt.Sprintf("Change: %s\nProject context:\n%v", changeName, contextInfo)

	response, err := o.provider.Chat(ctx, systemPrompt, fullPrompt)
	if err != nil {
		return nil, err
	}

	result := parseSDDPhaseResult(response)
	result.Phase = "sdd-apply"

	var action scanner.ProposedAction
	jsonErr := json.Unmarshal([]byte(extractJSON(response)), &action)
	if jsonErr == nil && action.Command != "" {
		action.Type = scanner.ActionCommand // Ensure type is set
		output, err := o.scanner.ValidateAndRunWithHitl(ctx, action, o.validator, "EVOLVE")
		if err != nil {
			return nil, err
		}
		result.DetailedReport = output
	}

	if result.Status == "" {
		result.Status = "success"
	}
	return result, nil
}

func parseSDDPhaseResult(raw string) *SDDPhaseResult {
	result := &SDDPhaseResult{Raw: raw, Status: "success", DetailedReport: raw}

	jsonBlock := extractJSON(raw)
	if strings.HasPrefix(strings.TrimSpace(jsonBlock), "{") {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(jsonBlock), &payload); err == nil {
			if status, ok := payload["status"].(string); ok {
				result.Status = status
			}
			if summary, ok := payload["executive_summary"].(string); ok {
				result.ExecutiveSummary = summary
			}
			if report, ok := payload["detailed_report"].(string); ok {
				result.DetailedReport = report
			}
			if next, ok := payload["next_recommended"].(string); ok {
				result.NextRecommended = next
			}
			if risks, ok := payload["risks"].(string); ok {
				result.Risks = risks
			}
			if resolution, ok := payload["skill_resolution"].(string); ok {
				result.SkillResolution = resolution
			}
			if artifacts, ok := payload["artifacts"].([]interface{}); ok {
				for _, artifact := range artifacts {
					if artifactStr, ok := artifact.(string); ok {
						result.Artifacts = append(result.Artifacts, artifactStr)
					}
				}
			}
		}
	}

	for _, line := range strings.Split(raw, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(strings.Trim(parts[0], "*` "))
		value := strings.TrimSpace(parts[1])
		switch strings.ToLower(key) {
		case "status", "**status**":
			result.Status = value
		case "executive_summary", "**executive_summary**", "summary", "**summary**":
			result.ExecutiveSummary = value
		case "detailed_report", "**detailed_report**":
			result.DetailedReport = value
		case "artifacts", "**artifacts**":
			for _, part := range strings.Split(value, "|") {
				cleaned := strings.TrimSpace(strings.Trim(part, "`[] "))
				if cleaned != "" {
					result.Artifacts = append(result.Artifacts, cleaned)
				}
			}
		case "next_recommended", "**next_recommended**", "next", "**next**":
			result.NextRecommended = value
		case "risks", "**risks**":
			result.Risks = value
		case "skill_resolution", "**skill_resolution**":
			result.SkillResolution = value
		}
	}

	if result.ExecutiveSummary == "" {
		result.ExecutiveSummary = strings.TrimSpace(strings.SplitN(raw, "\n", 2)[0])
	}

	return result
}

func normalizeSDDPhase(next string) string {
	next = strings.TrimSpace(strings.ToLower(next))
	if next == "" || next == "none" {
		return ""
	}

	known := []string{"sdd-init", "sdd-explore", "sdd-propose", "sdd-spec", "sdd-design", "sdd-tasks", "sdd-apply", "sdd-verify", "sdd-archive"}
	for _, phase := range known {
		if strings.Contains(next, phase) {
			return phase
		}
	}
	for _, phase := range known {
		if strings.Contains(next, strings.TrimPrefix(phase, "sdd-")) {
			return phase
		}
	}
	return next
}

func nextSDDPhase(current string, phases []string) string {
	for i, phase := range phases {
		if phase == current && i+1 < len(phases) {
			return phases[i+1]
		}
	}
	return ""
}

func determinePersistenceMode(content string) string {
	text := strings.ToLower(content)
	if strings.Contains(text, "hybrid") {
		return "hybrid"
	}
	if strings.Contains(text, "openspec") {
		return "openspec"
	}
	if strings.Contains(text, "none") {
		return "none"
	}
	return "engram"
}

func (o *Orchestrator) saveSDDArtifact(changeName, project, artifactKey, content, mode string) error {
	if mode == "none" {
		return nil
	}
	artifactID := fmt.Sprintf("sdd/%s/%s", changeName, artifactKey)
	if mode == "engram" || mode == "hybrid" {
		if err := o.memory.SaveFinding(artifactID, map[string]interface{}{
			"phase":   artifactKey,
			"project": project,
			"content": content,
		}); err != nil {
			return err
		}
	}
	if mode == "openspec" || mode == "hybrid" {
		dirPath := filepath.Join("openspec", "changes", changeName)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
		fileName := strings.ReplaceAll(artifactKey, "/", "-") + ".md"
		path := filepath.Join(dirPath, fileName)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func (o *Orchestrator) runSDDPhase(ctx context.Context, changeName, project, phase, input string) (*SDDPhaseResult, error) {
	skill, err := LoadSkill("Evolve", phase)
	if err != nil {
		return nil, err
	}
	common, err := LoadSkill("_shared", "sdd-phase-common")
	if err != nil {
		return nil, err
	}

	mode := determinePersistenceMode(skill.Content)
	contextInfo, _ := o.memory.RetrieveContext(48)
	summaryPrompt := fmt.Sprintf("Change: %s\nProject: %s\nPhase: %s\nPersistence Mode: %s\nPrevious context:\n%v\n\n%s\n\n%s", changeName, project, phase, mode, contextInfo, common.Content, skill.Content)
	response, err := o.provider.Chat(ctx, fmt.Sprintf("You are the %s SDD agent. Follow the shared SDD protocol and return a structured envelope for phase %s.", phase, phase), summaryPrompt)
	if err != nil {
		return nil, err
	}

	result := parseSDDPhaseResult(response)
	result.Phase = phase
	if result.Status == "" {
		result.Status = "success"
	}

	if err := o.saveSDDArtifact(changeName, project, phase, response, mode); err != nil {
		fmt.Printf("Warning: failed to persist artifact for %s: %v\n", phase, err)
	}

	if phase == "sdd-apply" {
		applyResult, err := o.runSubAgentApply(ctx, changeName)
		if err != nil {
			return nil, err
		}
		if applyResult.ExecutiveSummary != "" {
			result.ExecutiveSummary = applyResult.ExecutiveSummary
		}
		if applyResult.DetailedReport != "" {
			result.DetailedReport = applyResult.DetailedReport
		}
		result.Artifacts = append(result.Artifacts, applyResult.Artifacts...)
		result.NextRecommended = applyResult.NextRecommended
	}

	if len(result.Artifacts) == 0 {
		result.Artifacts = []string{fmt.Sprintf("sdd/%s/%s", changeName, phase)}
	}

	o.broadcastSDDProgress(result)
	return result, nil
}

func (o *Orchestrator) broadcastSDDProgress(result *SDDPhaseResult) {
	if o.dashboard == nil {
		return
	}
	o.dashboard.Broadcast(system.DashboardUpdate{
		Type:   "sdd_progress",
		Source: result.Phase,
		Phase:  result.Phase,
		Content: map[string]interface{}{
			"status":    result.Status,
			"summary":   result.ExecutiveSummary,
			"next":      result.NextRecommended,
			"artifacts": result.Artifacts,
		},
	})
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
		Type:        scanner.ActionCommand,
		Command:     fmt.Sprintf("sudo apt install -y %s", tool),
		Description: fmt.Sprintf("Install missing security tool: %s", tool),
	}
	o.scanner.ValidateAndRunWithHitl(ctx, action, o.validator, "IDENTIFY")
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
	fullPrompt := fmt.Sprintf("Change: %s\nProject context: %v\nTask: %s", skillName, contextInfo, input)

	response, err := o.provider.Chat(ctx, systemPrompt, fullPrompt)
	if err != nil {
		return nil, err
	}

	o.memory.SaveFinding(skillName, map[string]interface{}{"output": response})
	return &agents.Report{Summary: response}, nil
}

func (o *Orchestrator) findAlternativeTool(ctx context.Context, cat agents.NISTCategory, failedTool string, task string, err error) string {
	systemPrompt := fmt.Sprintf("You are a security tool advisor. A tool failed: %s with error: %v. Suggest an alternative tool for category %s and task: %s. Output only the tool name or 'NONE'.", failedTool, err, cat, task)
	response, chatErr := o.provider.Chat(ctx, systemPrompt, "")
	if chatErr != nil || strings.TrimSpace(response) == "NONE" {
		return ""
	}
	return strings.TrimSpace(response)
}

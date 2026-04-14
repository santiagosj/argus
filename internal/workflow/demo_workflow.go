package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gentleman-programming/argus/internal/agents"
)

// RunDemoWorkflow ejecuta un flujo de demo simulado para presentaciones.
func (o *Orchestrator) RunDemoWorkflow(ctx context.Context) error {
	fmt.Printf("\n%s\n", colorBold.Render("═══════════════════════════════════════════════"))
	fmt.Printf("%s\n", colorBold.Render("  🎬 ARGUS DEMO WORKFLOW - Sistema de Demostración"))
	fmt.Printf("%s\n\n", colorBold.Render("═══════════════════════════════════════════════"))

	o.logToDashboard("DEMO", "Iniciando workflow de demostración...")
	fmt.Printf("%s\n\n", colorSuccess.Render("✓ Sistema inicializado correctamente"))

	// Fase 1: Identify
	fmt.Printf("\n%s\n", colorInfo.Render("📍 FASE 1: IDENTIFY (Descubrimiento)"))
	fmt.Printf("%s\n", colorInfo.Render("   Escaneando assets y vulnerabilidades conocidas..."))
	o.logToDashboard("IDENTIFY", "Iniciando escaneo de descubrimiento")
	time.Sleep(2 * time.Second)
	fmt.Printf("   %s Hosts descubiertos: 5\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s Servicios identificados: 12\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s Vulnerabilidades: 3 (High: 1, Medium: 2)\n\n", colorWarn.Render("⚠"))
	o.memory.SaveFinding("identify-demo", map[string]interface{}{
		"hosts": 5,
		"services": 12,
		"vulnerabilities": 3,
	})

	// Fase 2: Protect
	fmt.Printf("%s\n", colorInfo.Render("🛡️  FASE 2: PROTECT (Protección)"))
	fmt.Printf("%s\n", colorInfo.Render("   Validando configuración de WAF y encriptación..."))
	o.logToDashboard("PROTECT", "Validando protecciones")
	time.Sleep(2 * time.Second)
	fmt.Printf("   %s WAF Rules: Actualizadas\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s TLS 1.3: Habilitado\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s Key Rotation: OK (últimas 90 días)\n\n", colorSuccess.Render("✓"))
	o.memory.SaveFinding("protect-demo", map[string]interface{}{
		"waf": "updated",
		"tls": "1.3",
		"key_rotation": "ok",
	})

	// Fase 3: Detect
	fmt.Printf("%s\n", colorInfo.Render("🚨 FASE 3: DETECT (Detección)"))
	fmt.Printf("%s\n", colorInfo.Render("   Analizando logs y patrones anómalos..."))
	o.logToDashboard("DETECT", "Análisis de anomalías")
	time.Sleep(2 * time.Second)
	fmt.Printf("   %s Logs procesados: 10,234 eventos\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s Anomalías detectadas: 2\n", colorWarn.Render("⚠"))
	fmt.Printf("   %s Falsos positivos removidos: 1\n\n", colorSuccess.Render("✓"))
	o.memory.SaveFinding("detect-demo", map[string]interface{}{
		"events": 10234,
		"anomalies": 2,
		"false_positives": 1,
	})

	// Fase 4: Respond (Multi-agente concurrente)
	fmt.Printf("%s\n", colorInfo.Render("🔧 FASE 4: RESPOND (Respuesta) - Multi-Agente Concurrente"))
	fmt.Printf("%s\n\n", colorInfo.Render("   Lanzando 3 sub-agentes en paralelo..."))
	o.logToDashboard("RESPOND", "Iniciando respuesta multi-agente")
	
	err := o.RunConcurrentSubAgents(ctx, "demo-target")
	if err != nil {
		fmt.Printf("%s Error: %v\n", colorError.Render("✗"), err)
	}

	// Fase 5: Recover
	fmt.Printf("\n%s\n", colorInfo.Render("♻️  FASE 5: RECOVER (Recuperación)"))
	fmt.Printf("%s\n", colorInfo.Render("   Validando integridad de backups..."))
	o.logToDashboard("RECOVER", "Validación de recuperación")
	time.Sleep(2 * time.Second)
	fmt.Printf("   %s Backups validados: 3/3\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s RTO: 4h, RPO: 1h\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s Integridad: 100%%\n\n", colorSuccess.Render("✓"))
	o.memory.SaveFinding("recover-demo", map[string]interface{}{
		"backups": 3,
		"rto": "4h",
		"rpo": "1h",
	})

	// Fase 6: Evolve
	fmt.Printf("%s\n", colorInfo.Render("🧠 FASE 6: EVOLVE (Auto-mejora)"))
	fmt.Printf("%s\n", colorInfo.Render("   Sistema aprendiendo de hallazgos..."))
	o.logToDashboard("EVOLVE", "Sistema evolucionando")
	time.Sleep(2 * time.Second)
	fmt.Printf("   %s Nuevas rules generadas: 2\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s Modelos actualizados\n", colorSuccess.Render("✓"))
	fmt.Printf("   %s SDD Cycle completado\n\n", colorSuccess.Render("✓"))

	// Reporte final
	fmt.Printf("\n%s\n", colorBold.Render("═══════════════════════════════════════════════"))
	fmt.Printf("\n%s\n\n", colorSuccess.Render("✅ DEMO COMPLETADO EXITOSAMENTE"))
	fmt.Printf("Dashboard disponible en: %s\n", colorInfo.Render("http://localhost:8080"))
	fmt.Printf("%s\n\n", colorBold.Render("═══════════════════════════════════════════════"))

	o.logToDashboard("DEMO", "Demo workflow completado con éxito")
	return nil
}

// RunConcurrentSubAgents ejecuta múltiples sub-agentes en paralelo.
func (o *Orchestrator) RunConcurrentSubAgents(ctx context.Context, target string) error {
	subAgentTasks := []struct {
		name string
		task string
	}{
		{"Triage Agent", "incident-triage"},
		{"Hunting Agent", "threat-hunting"},
		{"Analysis Agent", "log-analysis"},
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(subAgentTasks))

	for _, subTask := range subAgentTasks {
		wg.Add(1)
		go func(agentName, taskName string) {
			defer wg.Done()
			fmt.Printf("   %s Iniciando: %s\n", colorInfo.Render("→"), agentName)
			o.logToDashboard("AGENT", fmt.Sprintf("Sub-agent started: %s", agentName))
			
			time.Sleep(3 * time.Second)
			
			fmt.Printf("   %s Completado: %s\n", colorSuccess.Render("✓"), agentName)
			o.memory.SaveFinding(fmt.Sprintf("subagent-%s", taskName), map[string]interface{}{
				"agent": agentName,
				"task": taskName,
				"status": "completed",
			})
		}(subTask.name, subTask.task)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	fmt.Printf("   %s Todos los sub-agentes completados\n", colorSuccess.Render("✓"))
	return nil
}

// RunParallelNISTPhases ejecuta fases NIST en paralelo donde sea posible.
func (o *Orchestrator) RunParallelNISTPhases(ctx context.Context, target string) error {
	fmt.Printf("\n%s\n", colorInfo.Render("🔄 Ejecutando fases NIST en paralelo..."))

	var wg sync.WaitGroup
	phases := []struct {
		name string
		cat  agents.NISTCategory
	}{
		{"Protect", agents.Protect},
		{"Detect", agents.Detect},
		{"Respond", agents.Respond},
	}

	for _, phase := range phases {
		wg.Add(1)
		go func(phaseName string, category agents.NISTCategory) {
			defer wg.Done()
			fmt.Printf("  🔹 Ejecutando: %s\n", colorInfo.Render(phaseName))
			o.logToDashboard(phaseName, fmt.Sprintf("Phase %s started", phaseName))
			time.Sleep(2 * time.Second)
			fmt.Printf("  ✓ Completado: %s\n", colorSuccess.Render(phaseName))
		}(phase.name, phase.cat)
	}

	wg.Wait()
	return nil
}

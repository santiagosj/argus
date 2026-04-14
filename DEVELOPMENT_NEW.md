# Argus Development Guide (Updated)

## Current Status: v3.0.0-alpha - Demo Ready 🎬

### What's New in Phase 2 (Just Implemented)

✅ **Multi-Agent Concurrency**
- Sub-agents execute in parallel with `sync.WaitGroup`
- Isolated memory contexts per agent
- Full audit trail for each agent

✅ **Expanded Skills Library**  
- 8+ new skills for Protect, Detect, Respond, Recover phases
- Professional Markdown format with examples
- Auto-loaded discovery

✅ **Professional Demo Workflow**
- `argus demo` command shows all 6 NIST phases
- Pre-configured with simulated findings
- Real-time dashboard updates

✅ **Enhanced Dashboard**
- Tabbed interface by NIST phase
- Real-time filtering
- Timestamp tracking
- Entity view (IPs, CVEs, URLs, CWEs)

✅ **Production-Grade Documentation**
- README.md - Professional intro
- TESTING.md - Test scenarios
- DEMO_SCRIPT.md - Presentation guide
- DEVELOPMENT.md - This guide

✅ **Build & Deployment**
- Makefile for easy building
- Dockerfile for containerization
- install.sh for setup
- Cross-platform build targets

---

## Architecture Overview

```
┌─────────────────────────┐
│   CLI / Interactive TUI │
│  (demo, run, learn)     │
└──────────┬──────────────┘
           │
      ┌────▼─────────────────────┐
      │   Orchestrator           │
      │   └─ RunDemoWorkflow()   │
      │   └─ RunConcurrentSubs() │
      │   └─ Error Recovery      │
      └────┬────────┬──┬─────────┘
           │        │  │
    ┌──────▼──┐ ┌──▼──▼┐ ┌─────────┐
    │ Engram   │ │Skills │ Dashboard│
    │ (Memory) │ │(Tools)│ (Web)   │
    └──────────┘ └──────┘ └─────────┘
         │
    ┌────▼──────────────┐
    │  AI Providers    │
    │  - Ollama        │
    │  - Claude/OpenAI │
    └──────────────────┘
```

---

## Project Structure

```
argus/
├── cmd/argus/main.go                    # CLI entry
├── internal/
│   ├── app/app.go                       # App init + CLI logic
│   ├── agents/
│   │   ├── interface.go                 # Agent contract
│   │   └── ollama.go                    # LLM provider
│   ├── components/
│   │   ├── engram/memory.go             # Persistence + entities
│   │   ├── reports/generator.go         # Report builder
│   │   ├── scanner/executor.go          # Command runner
│   │   └── skills/generator.go          # Skill loader
│   ├── system/
│   │   ├── audit.go                     # Audit logging
│   │   ├── config.go                    # Config loader
│   │   └── server.go                    # Dashboard (SSE)
│   ├── tui/
│   │   ├── model.go                     # Bubble Tea state
│   │   ├── screens/                     # UI screens
│   │   └── styles/banner.go             # ASCII art
│   └── workflow/
│       ├── orchestrator.go              # Main orchestrator
│       ├── demo_workflow.go             # Demo + concurrency
│       └── skills.go                    # Skill management
├── skills/
│   ├── Identify/        (nmap, nuclei, recon)
│   ├── Protect/         (waf-validation, encryption-audit)
│   ├── Detect/          (log-analysis, ids-analysis)
│   ├── Respond/         (incident-triage, threat-hunting)
│   ├── Recover/         (backup-integrity, post-incident-hardening)
│   ├── Evolve/          (sdd phases)
│   └── _shared/         (common patterns)
├── config.yaml                          # Configuration
├── go.mod & go.sum                      # Dependencies
├── README.md                            # Main docs
├── README_DEMO.md                       # Quick start
├── TESTING.md                           # Test guide
├── DEMO_SCRIPT.md                       # Presentation
├── DEVELOPMENT.md                       # This file
├── Dockerfile                           # Container build
├── Makefile                             # Build targets
├── install.sh                           # Setup script
└── validate-demo.sh                     # Pre-demo check
```

---

## Building & Running

### Quick Build
```bash
make build
```

### Run Modes
```bash
./argus demo                    # Demo workflow (all 6 phases)
./argus run 127.0.0.1:8080      # Scan target
./argus audit 192.168.1.0/24    # Audit network
./argus learn                   # TUI interface
./argus status                  # Show config
```

### Docker
```bash
make docker                     # Build image
make docker-run                 # Run container
```

### Pre-Release Validation
```bash
chmod +x validate-demo.sh
./validate-demo.sh              # Check all systems go
```

---

## Key Code Components

### 1. Orchestrator (`orchestrator.go`)
**Main autonomous decision loop:**
```go
func (o *Orchestrator) RunAutonomousWorkflow(ctx context.Context, target string) error
  // - Loops 5 times
  // - Calls LLM to decide: TOOL | SUB_ORCHESTRATOR | FINISH
  // - Executes with error recovery
  // - Updates dashboard in real-time
```

**Demo workflow with all 6 phases:**
```go
func (o *Orchestrator) RunDemoWorkflow(ctx context.Context) error
  // - Runs Identify, Protect, Detect, Respond, Recover, Evolve
  // - Uses RunConcurrentSubAgents() for parallelization
  // - Populates memory with demo data
  // - Shows on dashboard instantly
```

**Parallel sub-agent execution:**
```go
func (o *Orchestrator) RunConcurrentSubAgents(ctx context.Context, target string) error
  // - Spawns 3+ sub-agents with sync.WaitGroup
  // - No blocking between agents
  // - Collects errors asynchronously
```

### 2. Engram Memory (`engram/memory.go`)
**Entity extraction:**
```go
func (s *SQLiteMemory) extractEntities(text string)
  // - Detects IPs, URLs, CVEs, CWEs
  // - Stores separately in entities table
  // - Enables fast searching
```

**Smart cleanup:**
```go
func (s *SQLiteMemory) Cleanup()
  // - Deletes findings older than TTL
  // - Removes orphaned entities
  // - Runs automatically
```

### 3. Dashboard (`system/server.go`)
**Server-Sent Events:**
```go
func (d *WebDashboard) Broadcast(update DashboardUpdate)
  // - Sends real-time events to all clients
  // - Dashboard auto-updates without polling
  // - Tabs filter by NIST phase
```

### 4. Skills System (`skills/` + `workflow/skills.go`)
**Format:** Markdown with Objective, Context, Instructions
**Discovery:** Auto-loaded by phase/tool name
**Extension:** Add new skills by creating `.md` files

---

## Demonstration Flow

### For Presentations (10 min)
```
1. ./argus demo
2. Open http://localhost:8080
3. Show dashboard tabs updating real-time
4. Explain multi-agent concurrency
5. Query database: sqlite3 argus_memory.db
6. Show audit log: tail -f argus_audit.jsonl
```

Full script: See [DEMO_SCRIPT.md](DEMO_SCRIPT.md)

---

## Testing Checklist

Before demo/release:
```bash
✓ make build                    # Compiles
✓ ./argus help                  # Shows commands
✓ ./argus status                # Shows config
✓ ./argus demo                  # Runs without error
✓ curl localhost:8080           # Dashboard responds
✓ sqlite3 argus_memory.db       # Database has data
✓ validate-demo.sh              # All systems green
```

---

## Code Style & Conventions

### Go Packages
- `agents/` - AI provider adapters
- `components/` - Feature modules (Engram, Skills, Reports)
- `system/` - Utilities (Config, Audit, Server)
- `tui/` - User interface components
- `workflow/` - Orchestration logic

### Color Coding (Feedback)
```go
colorSuccess = "#55ff55"  // ✓ Completed
colorInfo    = "#00d4ff"  // → Progress
colorWarn    = "#ffaa00"  // ⚠ Warnings
colorError   = "#ff5555"  // ✗ Errors
```

### Error Handling
```go
// Always log to both console AND dashboard
o.logToDashboard("ERROR", fmt.Sprintf("Failed: %v", err))
fmt.Printf(colorError.Render("✗ Error: %v\n"), err)
```

---

## Next Development Priorities

### Short Term (Next Release)
- [ ] Claude/OpenAI provider adapters
- [ ] Workflow templates (save/load common patterns)
- [ ] Auto skill generation for new tools
- [ ] TUI improvements (drag-drop skills)

### Medium Term
- [ ] Slack/Teams integration
- [ ] Jira ticket creation
- [ ] Webhook triggers
- [ ] PDF reports

### Long Term
- [ ] Contextual learning (skill effectiveness scoring)
- [ ] RBAC & multi-user approvals
- [ ] Distributed workflow across machines
- [ ] Commercial SaaS platform

---

## Contributing

1. Create feature branch: `git checkout -b feature/xyz`
2. Follow code style above
3. Add test coverage where possible
4. Update README/docs if needed
5. Submit PR with clear description

---

## Support

- **Issues:** File in GitHub
- **Discussions:** GitHub Discussions board
- **Security:** Report privately to [security@url]

---

**Last Updated:** April 9, 2026 | **Status:** v3.0.0-alpha (Production-Ready Demo)

# 🔍 Argus - Cognitive Cyber-Security Framework

**The "Copilot" for Security Engineers**

Argus es un orquestador autónomo de ciberseguridad que automatiza pentesting y auditorías de seguridad mediante agentes IA con memoria persistente, todo alineado al **NIST Cybersecurity Framework**.

---

## 🚀 Quick Start

### Instalación
```bash
# Prerequisitos: Go 1.25+, Ollama (o LM Studio para modelos locales)
git clone https://github.com/gentleman-programming/argus
cd argus
go build ./cmd/argus -o argus
```

### Primeros pasos
```bash
# Ver ayuda
./argus help

# Ver demo (sin requerir target específico)
./argus demo

# Ejecutar sobre un target real
./argus run localhost:8080

# Abrir menú interactivo de learning
./argus learn

# Ver status
./argus status
```

### Ver Dashboard en Tiempo Real
Una vez ejecutando un workflow:
```
http://localhost:8080
```
El dashboard mostrará hallazgos en tiempo real, categorizados por fase NIST.

---

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────┐
│        CLI / Interactive TUI            │
│   (Target Selection, Skill Selector)    │
└────────────┬────────────────────────────┘
             │
        ┌────▼──────────────────────────┐
        │   Orchestrator (Brain)        │
        │ - Decision Loop               │
        │ - Multi-Agent Spawner         │
        │ - Error Recovery              │
        └────┬──────────┬───────────────┘
             │          └──────────────────┐
      ┌──────▼────────┐            ┌──────▼────────┐
      │ Engram        │            │ Skills        │
      │ (Persistent   │            │ (NIST Lib)    │
      │  Memory)      │            └───────────────┘
      └───────────────┘
             │
      ┌──────▼─────────────────────┐
      │ Agent Providers            │
      │ - Ollama (Local)           │
      │ - Claude/OpenAI (Cloud)    │
      └────────────────────────────┘
```

### Componentes Clave

| Componente | Descripción |
|-----------|-------------|
| **Orchestrator** | Motor de decisiones autónomo que planifica y ejecuta workflows NIST |
| **Engram** | Memoria persistente (SQLite) con extracción de entidades (IPs, CVEs, URLs) |
| **Skills** | Biblioteca de prompts para herramientas de seguridad (Identify, Protect, Detect, Respond, Recover, Evolve) |
| **Dashboard** | UI web en tiempo real con Server-Sent Events (SSE) |
| **Agents** | Adaptadores modulares para Ollama, Claude, OpenAI |
| **HITL** | Human-in-the-Loop validation para operaciones críticas |

---

## 🎯 Flujo NIST Integrado

Argus ejecuta flujos alineados a las 5 funciones del NIST CSF:

1. **Identify** 🔍
   - Recon, NMAP, Nuclei
   - Descubrimiento de activos
   - Mapeo de vulnerabilidades conocidas

2. **Protect** 🛡️
   - Hardening checks
   - WAF rules validation
   - Compliance verification

3. **Detect** 🚨
   - Anomaly detection
   - Log analysis
   - Traffic inspection

4. **Respond** 🔧
   - Incident triage
   - Automated responses
   - Root cause analysis

5. **Recover** ♻️
   - Restoration validation
   - Backup verification
   - Post-incident hardening

6. **Evolve** 🧠
   - Self-improvement via SDD
   - Auto-skill generation
   - Pattern learning

---

## 💡 Características Principales

### ✅ Cognición Local & Privada
- Soporte para Ollama (modelos locales sin enviar datos a la nube)
- Alternativa: Claude/OpenAI para análisis en la nube

### ✅ Anti-Amnesia (Engram)
- Cada hallazgo se persiste en SQLite
- Extracción automática de entidades (IPs, CVEs, URLs, CWEs)
- Deduplicación de findings
- TTL configurable para limpiar hallazgos antiguos

### ✅ Multi-Agente Concurrente
- Orquestación de sub-agentes en paralelo
- Aislamiento de memoria por sub-agente
- Comunicación asincrónica

### ✅ Human-in-the-Loop (HITL)
- Validación de acciones antes de ejecución
- Logs auditables de todas las decisiones
- Modo dry-run para preview

### ✅ Recuperación Automática de Errores
- Si una herramienta falla, propone alternativa
- Reintentos inteligentes
- Logging completo de fallos

### ✅ Dashboard Web Profesional
- Visualización en tiempo real de events
- Tabs separadas por fase NIST
- Tablas filtrables de hallazgos
- Timestamps y source tracking

---

## 📋 Ejemplos de Uso

### Ejemplo 1: Scan Rápido
```bash
./argus run 192.168.1.1
```
Ejecuta Identify + Protect check automáticamente.

### Ejemplo 2: Auditoría Completa
```bash
./argus run --full 10.0.0.0/24
```
Ejecuta todas las fases NIST en paralelo.

### Ejemplo 3: Evolve (Auto-mejora)
```bash
./argus run evolve
```
Sistema se auto-mejora mediante Spec-Driven Development.

---

## ⚙️ Configuración

Editar `config.yaml`:

```yaml
ai:
  provider: "ollama"           # o "claude", "openai"
  model: "mistral:latest"      # modelo a usar
  base_url: "http://localhost:11434"

persistence:
  type: "sqlite"               # tipo de persistencia
  path: "argus_memory.db"      # ruta de memoria

tools:
  auto_install: false          # instalar herramientas auto
```

---

## 🧪 Demo Pre-configurado

```bash
./argus demo
```

Ejecuta un workflow de demo que:
- ✅ Muestra scanning simulado
- ✅ Demuestra findings en dashboard
- ✅ Muestra multi-agente en acción
- ✅ Genera reporte final
- ✅ Toma ~2-3 minutos

---

## 📊 Roadmap

- [ ] **Fase 2:** Claude/OpenAI adapters, workflow templates, skill auto-generation
- [ ] **Fase 3:** Contextual learning, integration webhooks (Slack, Jira)
- [ ] **Fase 4:** RBAC, multi-user approvals, PDF reports

---

## 🔐 Seguridad & Privacidad

- Todo procesamiento es local por defecto
- Soporte para modelos privados vía Ollama
- Logs auditables en `argus_audit.jsonl`
- Sin envío de datos sensibles a servicios externos (por defecto)

---

## 📝 Licencia

MIT - Gentleman Programming

---

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Por favor:
1. Fork el repo
2. Crea rama para tu feature
3. Envía PR

---

## 📧 Contacto

Para preguntas, sugerencias o reportar bugs:
- Issues: GitHub Issues
- Email: [tu-email]

---

**Argus v3.0.0-alpha** | Made with ❤️ by Gentleman Programming

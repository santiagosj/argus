# 🔍 Argus - Cognitive Cyber-Security Framework

**The "Copilot" for Security Engineers**

Argus es un orquestador autónomo de ciberseguridad que automatiza pentesting y auditorías de seguridad mediante agentes IA con memoria persistente, todo alineado al **NIST Cybersecurity Framework**.

![Latest Release](https://img.shields.io/badge/version-v3.0.0--alpha-blue)
![License](https://img.shields.io/badge/license-MIT-green)

---

## 🚀 Quick Start

### Instalación Rápida

```bash
# Prerequisitos: Go 1.25+
git clone https://github.com/gentleman-programming/argus
cd argus
go build ./cmd/argus -o argus
```

### Primeros pasos (sin dependencias)

```bash
# Demo preconfigurada
./argus demo

# Ver ayuda
./argus help

# Ejecutar sobre un target
./argus run localhost:8080
```

### Dashboard en Tiempo Real

Mientras ejecutas un workflow, abre en el navegador:
```
http://localhost:8080
```

---

## 📊 ¿Qué puede hacer Argus?

### ✅ Automatiza Pentesting Completo
- 🔍 Descubrimiento de activos (Identify)
- 🛡️  Validación de protecciones (Protect)  
- 🚨 Detección de anomalías (Detect)
- 🔧 Respuesta a incidentes (Respond)
- ♻️  Recuperación de fallos (Recover)
- 🧠 Auto-mejora continua (Evolve)

### ✅ Cognición Local & Privada
- Modelos locales vía Ollama (sin enviar datos a cloud)
- API en nube (Claude/OpenAI) como alternativa
- Control total de tus datos

### ✅ Memoria Persistente (Engram)
- SQLite con extracción automática de entidades
- Deduplicación inteligente
- TTL configurable
- Búsqueda por similaridad

### ✅ Multi-Agente Concurrente
- Sub-agentes especializados
- Ejecución en paralelo
- Aislamiento de contexto
- Comunicación asincrónica

### ✅ Human-in-the-Loop
- Validación de acciones antes de ejecutar
- Logs auditables completos
- Modo dry-run para preview

---

## 🏗️ Arquitectura

```
┌──────────────────────────────┐
│   CLI / Interactive TUI      │
└────────┬─────────────────────┘
         │
    ┌────▼──────────────────┐
    │  Orchestrator (Brain) │
    │  - Decision Loop      │
    │  - Multi-Agent Spawn  │
    │  - Error Recovery     │
    └────┬─────────────────┘
         │
    ┌────┴─────┬──────────┬──────────┐
    │ Engram    │ Skills   │Dashboard │
    │(Memory)   │(Tools)   │(Web UI)  │
    └───────────┴──────────┴──────────┘
         │
    ┌────▼──────────────────────┐
    │  AI Providers            │
    │  - Ollama (Local)        │
    │  - Claude/OpenAI (Cloud) │
    └──────────────────────────┘
```

---

## 📚 Componentes

| Componente | Descripción |
|-----------|-------------|
| **Orchestrator** | Motor autónomo de decisiones |
| **Engram** | Memoria persistente inteligente |
| **Skills** | Prompts para herramientas de seguridad |
| **Dashboard** | UI web con SSE en tiempo real |
| **Agents** | Adaptadores para IA (Ollama, Claude, OpenAI) |
| **HITL** | Validación human-in-the-loop |

---

## 🎯 Flujo NIST Integrado

Argus ejecuta automáticamente las 6 funciones del NIST Cybersecurity Framework:

**1. Identify 🔍** - Descubre activos y vulnerabilidades  
**2. Protect 🛡️** - Valida configuraciones de seguridad  
**3. Detect 🚨** - Analiza logs y detecta anomalías  
**4. Respond 🔧** - Clasifica y responde a incidentes  
**5. Recover ♻️** - Verifica recuperación y backups  
**6. Evolve 🧠** - Auto-mejora mediante SDD  

---

## 📖 Ejemplos de Uso

### Ejemplo 1: Demo (sin dependencias)
```bash
./argus demo
```
Flujo completo demostrativo con todas las fases.

### Ejemplo 2: Scan Rápido
```bash
./argus run 192.168.1.1
```
Ejecuta Identify + Protect automáticamente.

### Ejemplo 3: Auditoría Completa
```bash
./argus run --full 10.0.0.0/24
```
Todas las fases NIST en paralelo.

### Ejemplo 4: Learning Mode
```bash
./argus learn
```
Abre interfaz interactiva para explorar herramientas.

---

## ⚙️ Configuración

Edita `config.yaml`:

```yaml
ai:
  provider: "ollama"           # o "claude", "openai"
  model: "mistral:latest"
  base_url: "http://localhost:11434"

persistence:
  type: "sqlite"
  path: "argus_memory.db"

tools:
  auto_install: false
```

---

## 🧪 Testing

Para validar la instalación:

```bash
./argus help               # ✓ Muestra commands
./argus status             # ✓ Muestra config
./argus demo               # ✓ Ejecuta demo
curl localhost:8080        # ✓ Dashboard activo
```

Ver [TESTING.md](TESTING.md) para más scenarios.

---

## 🛡️ Seguridad & Privacidad

- ✅ Ejecución local por defecto
- ✅ Sin envío de datos a servidores externos (por defecto)
- ✅ Logs auditables en `argus_audit.jsonl`
- ✅ Validación human-in-the-loop
- ✅ Control total de datos sensibles

---

## 🗺️ Roadmap

- [x] Core orchestrator con decisiones autónomas
- [x] Multi-agent concurrente
- [x] Engram con extracción de entidades
- [x] Dashboard web time-real
- [x] NIST phases mapping
- [ ] Claude/OpenAI adapters
- [ ] Workflow templates
- [ ] Auto skill generation
- [ ] Contextual learning
- [ ] RBAC & multi-user
- [ ] PDF reports

---

## 📝 Licencia

MIT 

---

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Por favor:
1. Fork el repo
2. Crea rama `feature/tu-feature`
3. Envía pull request

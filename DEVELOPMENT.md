# Argus: Cognitive Cyber-Security Framework (Development Guide)

## 🎯 Visión del Proyecto
`Argus` no es una herramienta de escaneo estática; es un **Orquestador Cognitivo** diseñado para ingenieros de ciberseguridad. Su objetivo es actuar como un "copiloto" que gestiona sub-agentes granulares, manteniendo contextos pequeños para optimizar el uso de tokens y evitar la amnesia del modelo, todo mientras se ejecuta eficientemente en hardware doméstico (GPU decente + Ollama/Local).

## 🏗️ Arquitectura (Inspirada en Gentle-AI)
El proyecto sigue un patrón de **Desacoplamiento Total**:

1.  **Adapters (`internal/agents`):** Abstraen el proveedor de IA (Ollama para local, Claude/OpenAI para nube).
2.  **Components (`internal/components`):** "Supercargadores" que se inyectan en los agentes.
    *   **Engram:** Memoria persistente para evitar repetir hallazgos.
    *   **NIST:** Mapeo de lógica al framework NIST CSF.
3.  **Orquestador (`internal/workflow`):** El cerebro que decide qué sub-agente disparar, qué skill cargar y cómo pasar la estafeta entre agentes.
4.  **TUI (`internal/tui`):** Interfaz rica basada en Bubble Tea para que el humano valide y dirija el flujo.

## 🧠 Flujo de Trabajo Cognitivo
1.  **Selección NIST:** El humano elige una fase (ej. Identify).
2.  **Carga de Skills:** El orquestador busca en `skills/` el Markdown correspondiente.
3.  **Sub-Agente Quirúrgico:** Se lanza un agente con un system prompt mínimo + hallazgos previos de Engram.
4.  **Validación Humana:** El agente propone, el humano autoriza.

---

## 🚀 Próximos Pasos (Roadmap de Implementación)

### 1. Sistema de Auto-Evolución (Self-Development via SDD) 🛠️
Implementar la capacidad del framework para desarrollarse a sí mismo:
*   **Generador de Skills:** Un agente especializado que lea la documentación de una nueva herramienta de seguridad (ej. `zap-cli`) y genere automáticamente un archivo `SKILL.md` en la categoría NIST correspondiente.
*   **Agent-Factory via SDD:** Utilizar Spec-Driven Development para que el framework, ante una tarea desconocida, diseñe la especificación de un nuevo sub-agente, genere su lógica de adapter y lo integre al orquestador.

### 2. Validación Humana & "Safety Guards" 🛡️
*   Crear una pantalla de "Propuesta de Acción" donde el agente explique: "Voy a ejecutar `nmap -sV target`. ¿Autorizas?".
*   Implementar un sandbox para la ejecución de comandos generados por la IA.

### 3. Generador de Reportes Cognitivos 📝
*   Componente que consolide los hallazgos de múltiples sub-agentes en un reporte ejecutivo y técnico (Markdown/PDF) alineado con los controles del NIST.

### 4. Integración Profunda con Engram 💾
*   Refinar la recuperación de contexto para que los sub-agentes no solo lean "texto", sino "entidades de seguridad" (IPs, Vulnerabilidades, Credenciales) detectadas previamente.

---

## 📖 Nota para Sesiones Futuras
Cuando leas este proyecto, enfócate en la **interfaz `CyberAgent`** y el **`Orchestrator`**. Cualquier nueva funcionalidad debe ser un `Component` o una `Skill`. Para la auto-generación de código, utiliza el flujo SDD (Explore -> Propose -> Spec -> Implement) para asegurar que el framework crezca de forma estructurada y no caótica.

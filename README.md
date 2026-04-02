# Argus - Cognitive Cyber-Security Framework

**Un orquestador cognitivo de ciberseguridad inspirado en el Gentleman Stack.**

`Argus` es un framework CLI que actúa como el "piloto" del ingeniero de seguridad. Basado en el **NIST Cybersecurity Framework (CSF)**, automatiza tareas repetitivas de pentesting y validación mediante agentes y sub-agentes con memoria persistente.

## Pilares de Argus

1. **Cognición Local & Nube:** Soporte para modelos locales (vía Ollama/LM Studio) para privacidad en datos sensibles, o APIs en la nube para análisis masivo.
2. **Anti-Amnesia (Engram):** Cada hallazgo, vulnerabilidad o configuración se persiste en la memoria del agente para evitar el consumo excesivo de tokens en re-escaneos.
3. **Flujo NIST (Identify, Protect, Detect, Respond, Recover):** Las "Skills" del framework se mapean directamente a estas categorías.
4. **Agentes Granulares:** En lugar de un solo agente con un contexto enorme, `Argus` lanza sub-agentes específicos (ej. `recon-agent`, `exploit-agent`, `report-agent`) con skills mínimas necesarias.
5. **Human-in-the-Loop:** El ingeniero siempre tiene el control, validando las propuestas del agente antes de la ejecución de comandos destructivos.

## Arquitectura (Basada en Gentle-AI)

- `internal/agents`: Adapters para diferentes backends de IA (Ollama, Claude, OpenAI).
- `internal/components`: Supercargadores para el agente (Engram, NIST Skills, Report Generator).
- `internal/tui`: Interfaz de usuario rica (Bubble Tea) para el dashboard de ciberseguridad.
- `skills/`: Biblioteca de prompts y scripts especializados por categoría NIST.

## Referencias
- Framework: [Cybersecurity Framework (NIST)](https://cybersecurityframework.io/)
- Core: Gentleman Programming (Gentle-AI, Engram, Agent-Teams-Lite)

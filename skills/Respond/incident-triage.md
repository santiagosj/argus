# Respond: Incident Triage Skill

## Objetivo
Clasificar y priorizar incidentes de seguridad para respuesta efectiva.

## Contexto
Ante un incidente, la rapidez y precisión de triage son críticas. Esta skill:
- Valida que es un verdadero positivo
- Clasifica tipo
- Asigna severidad
- Identifica sistemas afectados

## Instrucciones

1. **Validar incidente**
   - Confirmar verdadero positivo vs false positive
   - Recolectar evidencia inicial
   - Documentar timeline

2. **Clasificación**
   - Tipo de incidente (breach, malware, DDoS, defacement, etc.)
   - Scope (single system vs enterprise)
   - Intent (accidental, targeted, opportunistic)

3. **Asignación de severidad**
   - CVSS scoring si aplica
   - Business impact assessment
   - Data sensitivity

4. **Escalation decision**
   - Quién debe notificarse
   - Nivel de urgencia
   - Recursos necesarios

## Ejemplo Output
```json
{
  "incident_type": "credential_compromise",
  "severity": "Critical",
  "affected_systems": ["web-01", "db-primary"],
  "escalate_to": ["CISO", "IR-team"],
  "recommendation": "Initiate incident response playbook"
}
```

# Recover: Post-Incident Hardening Skill

## Objetivo
Endurecer el sistema después de un incidente para prevenir recurrencia.

## Contexto
La recuperación no es solo volver online, es mejorar defensas. Esta skill:
- Identifica vulnerabilidades explotadas
- Recomienda patches
- Configura mitigaciones
- Documenta lecciones aprendidas

## Instrucciones

1. **Analizar incidente**
   - Root cause analysis
   - Systems affected
   - Attack vector

2. **Identificar gaps de seguridad**
   - Vulnerabilidades explotadas
   - Controles fallidos
   - Configuration issues

3. **Aplicar hardening**
   - Security patches
   - Configuration changes
   - Network segmentation
   - Access control updates

4. **Output format**
   - Cambios realizados
   - Before/after comparison
   - Verification steps

## Ejemplo Output
```json
{
  "root_cause": "unpatched RCE vulnerability",
  "hardening_applied": [
    "Applied patch 2026-04-05",
    "Enabled WAF rule for CVE-2026-1234",
    "Implemented network segmentation"
  ],
  "verification": "passed",
  "recommendation": "Schedule security training"
}
```

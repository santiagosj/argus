# Respond: Threat Hunting Skill

## Objetivo
Búsqueda proactiva de amenazas y IOCs (Indicators of Compromise) en la infraestructura.

## Contexto
La caza de amenazas es actividad proactiva. Esta skill:
- Define hipótesis de amenaza
- Busca indicators
- Valida presencia o ausencia
- Documenta findings

## Instrucciones

1. **Formular hipótesis de amenaza**
   - Basada en threat intelligence
   - MITRE ATT&CK techniques
   - Known attacker patterns

2. **Definir IOCs a buscar**
   - Hash de malware conocido
   - Dominio malicioso
   - IP address
   - Registry keys
   - Scheduled tasks anómalas

3. **Ejecutar búsqueda**
   - Query EDR/SIEM data
   - Verificar múltiples fuentes
   - Correlate findings

4. **Output format**
   - IOCs encontrados
   - Timeline de actividad
   - Recomendaciones

## Ejemplo Output
```json
{
  "hypothesis": "Lazarus-like lateral movement",
  "iocs_found": ["cmd.exe with /c rundll32"],
  "affected_hosts": ["workstation-23", "server-app-02"],
  "recommendation": "Isolate hosts and begin forensic analysis"
}
```

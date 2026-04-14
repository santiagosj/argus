# Detect: IDS/IPS Analysis Skill

## Objetivo
Analizar alertas de Intrusion Detection/Prevention Systems para identificar intentos de intrusión.

## Contexto
IDS/IPS proveen visibilidad de la red. Esta skill:
- Correlaciona alertas de múltiples sensores
- Identifica false positives
- Prioriza amenazas reales
- Mapea al MITRE ATT&CK framework

## Instrucciones

1. **Recolectar alertas IDS/IPS**
   - Últimas 24 horas
   - Agrupar por tipo
   - Filtrar por severidad

2. **Análisis de amenazas**
   - Source IP reputation
   - Payload analysis
   - Protocol anomalies
   - Port scanning patterns

3. **Correlación de eventos**
   - Multiple alert triggers = coordinated attack
   - Timeline reconstruction
   - Attack chain identification

4. **Output format**
   - Severity-sorted alerts
   - Attack classification
   - Recommended actions

## Ejemplo Output
```json
{
  "total_alerts": 23,
  "critical": 2,
  "suspicious_ips": ["203.x.x.x", "185.x.x.x"],
  "attack_patterns": ["port scanning", "payload injection"]
}
```

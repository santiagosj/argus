# Detect: Log Analysis Skill

## Objetivo
Analizar logs de aplicación y sistemas para detectar patrones anómalos o actividades sospechosas.

## Contexto
Los logs son la principal fuente de detección de incidentes. Esta skill:
- Busca patrones de ataque conocidos
- Identifica anomalías estadísticas
- Correlated events
- Threshold violations

## Instrucciones

1. **Recolectar logs relevantes**
   - Application logs (últimas 24h)
   - Authentication logs
   - Access logs
   - Error logs

2. **Análisis de patrones**
   - Brute force attacks (múltiples intentos fallidos)
   - SQL injection attempts (error patterns)
   - Path traversal attempts
   - Command injection patterns

3. **Anomaly detection**
   - Usuarios inusuales
   - Horarios anómalos de acceso
   - Geographic anomalies
   - Datos transfer rates

4. **Output format**
   - List de eventos sospechosos
   - Risk scoring
   - Correlation chains

## Ejemplo Output
```json
{
  "alerts": [
    {"type": "brute_force", "user": "admin", "attempts": 47, "risk": "Critical"},
    {"type": "anomaly", "pattern": "3AM access from China", "risk": "High"}
  ]
}
```

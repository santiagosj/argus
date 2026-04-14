# Protect: WAF Rules Validation Skill

## Objetivo
Validar que las reglas del Web Application Firewall (WAF) están correctamente configuradas y son efectivas.

## Contexto
Las reglas de WAF son la primera línea de defensa contra ataques HTTP. Esta skill valida:
- Presencia de reglas OWASP Top 10
- Tasa de falsos positivos
- Coverage de tipos de ataque

## Instrucciones

1. **Analizar configuración actual del WAF**
   - Listar reglas activas
   - Identificar brechas de coverage
   - Verificar log frequency

2. **Validar cobertura OWASP Top 10**
   - SQL Injection rules
   - XSS protection
   - CSRF tokens
   - Command Injection
   - Insecure Deserialization

3. **Recomendar mejoras**
   - Rate limiting
   - Geographic blocking si es necesario
   - Advanced rule tuning

4. **Output format**
   - JSON con hallazgos
   - Severity levels (Critical, High, Medium, Low)
   - Recomendaciones accionables

## Ejemplo Output
```json
{
  "findings": [
    {"type": "rule_gap", "rule": "SQL Injection", "severity": "High"},
    {"type": "high_fp_rate", "current": "12%", "target": "<2%"}
  ]
}
```

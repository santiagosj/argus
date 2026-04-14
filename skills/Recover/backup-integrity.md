# Recover: Backup Integrity Skill

## Objetivo
Validar que los backups son recuperables y no han sido comprometidos.

## Contexto
Los backups son el último recurso ante catastrophe. Esta skill:
- Valida integridad
- Verifica recoverability
- Testea datos
- Documenta RTO/RPO

## Instrucciones

1. **Verificar backup integrity**
   - Validar checksums
   - Verificar no corruption
   - Test restore proceso

2. **Verify backup availability**
   - Ubicaciones diversas
   - Offline backups
   - Encryption keys accessible

3. **RTO/RPO validation**
   - Recovery Time Objective
   - Recovery Point Objective
   - Meets SLAs

4. **Output format**
   - Status de cada backup
   - Integridad score
   - Last tested date

## Ejemplo Output
```json
{
  "backups": [
    {"location": "primary-vault", "status": "healthy", "tested": "2026-04-01"},
    {"location": "cold-storage", "status": "healthy"}
  ],
  "rto": "4h",
  "rpo": "1h",
  "recommendation": "All systems ready for recovery"
}
```

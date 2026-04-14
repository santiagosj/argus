# Protect: Encryption Audit Skill

## Objetivo
Auditar que todos los datos sensibles están siendo encriptados correctamente (transit + at rest).

## Contexto
El cifrado es un pilar fundamental de protección. Esta skill valida:
- TLS/SSL configuration
- Cipher suite strength
- Key rotation policies
- Database encryption

## Instrucciones

1. **Verificar Encryption in Transit**
   - Validar TLS 1.2+ en todos los endpoints
   - Listar cipher suites activos
   - Verificar certificate validity
   - Validar HSTS headers

2. **Verificar Encryption at Rest**
   - Database encryption status
   - File system encryption
   - Backup encryption
   - Key management system

3. **Key Rotation Audit**
   - Último cambio de keys
   - Frecuencia de rotación
   - Backup of old keys

4. **Output format**
   - JSON con status de cada componente
   - Severity de cualquier issue
   - Compliance map (PCI-DSS, HIPAA)

## Ejemplo Output
```json
{
  "tls_status": "TLS 1.3 enabled",
  "weak_ciphers": [],
  "key_rotation_days": 90,
  "last_rotation": "2026-03-15"
}
```

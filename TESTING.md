# 🧪 Testing & Demo Guide for Argus

## Quick Start sin Dependencies

Si no tienes Go instalado localmente, puedes:

1. **Usar el ejecutable**
   ```bash
   # En la demo, mostrar solo usando pre-compilado
   ./argus demo
   ```

2. **O compilar con Docker**
   ```bash
   docker build -t argus .
   docker run argus demo
   ```

## Testing Scenarios

### Scenario 1: Demo Rápida (2 minutos)
```bash
./argus demo
```
Muestra:
- ✅ Descubrimiento de assets
- ✅ Validación de protecciones
- ✅ Detección de anomalías
- ✅ Respuesta multi-agente (concurrente)
- ✅ Recuperación
- ✅ Auto-mejora

### Scenario 2: Workflow Real
```bash
# Requiere Ollama corriendo
ollama serve &
./argus run 127.0.0.1:8080
```

### Scenario 3: Dashboard en Vivo
Accede a http://localhost:8080 mientras ejecutas workflow

## Architecture Testing

### Test 1: CLI Interface
```bash
./argus help          # Should show commands
./argus status        # Should show config
./argus version       # Should show version
```

### Test 2: Dashboard SSE
```bash
curl http://localhost:8080/events
# Should receive server-sent events in real-time
```

### Test 3: Multi-Agent Concurrency
- El demo mostrará 3 sub-agentes ejecutando en paralelo
- Verifica que aparezcan timestamps paralelos en el dashboard

### Test 4: Memory Persistence
```bash
# Revisa argus_memory.db después de una ejecución
sqlite3 argus_memory.db "SELECT * FROM findings LIMIT 5;"
```

## Validación de Features

| Feature | Test | Expected |
|---------|------|----------|
| CLI | `./argus help` | Muestra todos comandos |
| Demo | `./argus demo` | Completa en ~2-3 min |
| Dashboard | `localhost:8080` | Muestra hallazgos en tiempo real |
| Multi-Agent | Demo workflow | 3+ sub-agentes en paralelo |
| Engram | Query DB | Activos con TTL |
| Entity Extraction | DB check | IPs, CVEs, URLs extraidos |
| Error Recovery | Fallo simulado | Propone alternativa |

## Demo Flow para Presentación

```
1. Abrir dashboard en navegador (localhost:8080)
2. Terminal 1: ./argus demo
3. Ver datos en tiempo real en dashboard
4. Mostrar tabs de fases NIST
5. Explicar multi-agente concurrente
6. Mostrar logs auditables (argus_audit.jsonl)
7. Revisar memoria (argus_memory.db)
```

## Troubleshooting

### "command not found: argus"
```bash
cd /path/to/argus
go build ./cmd/argus -o argus
./argus demo
```

### "Failed to load config"
```bash
# Usa defaults o crea config.yaml
cp config.yaml.example config.yaml
```

### Dashboard no muestra datos
- Verifica que está ejecutándose: `lsof -i :8080`
- Actualiza el navegador
- Abre DevTools console para ver errores

### Ollama connection failed
- Verifica: `ollama serve` en otra terminal
- O usa `--local` flag (cuando esté implementado)


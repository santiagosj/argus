# 🎬 Demo Script para Presentación

## Pre-presnetación (Setup)

### 1. Preparar ambiente (5 min antes)
```bash
# Terminal 1: Compilar
cd /home/kali/proyectos/argus
make build

# Terminal 2: Dashboard (o verlo en navegador durante demo)
# Dejarlo listo para http://localhost:8080

# Terminal 3: Verificar que todo funciona
./argus status
./argus help
```

### 2. Abrir navegador
Ir a: `http://localhost:8080` (estará vacío hasta que inicie workflow)

---

## 📊 Guión de Presentación (10 minutos)

### Slide 1: Introducción (1 min)
```
"Argus es un orquestador autónomo de ciberseguridad que actúa como copiloto para ingenieros de seguridad.

Características principales:
✓ Automatización de pentesting completo
✓ Cognición local & privada (sin enviar datos a la nube)
✓ Multi-agente concurrente
✓ Memoria persistente inteligente
✓ Flujo NIST integrado
✓ Dashboard en tiempo real
```

### Slide 2: Arquitectura (1 min)
Mostrar arquitectura en README.md

### Slide 3: Live Demo - Initializar (0.5 min)
```bash
# Terminal: El presentador ejecuta
./argus
```

Explicar qué pasa:
- "Argus arranca la TUI interactiva y permite seleccionar el target y la fase NIST desde la interfaz"
- "Veremos datos apareciendo en tiempo real en el dashboard"

### Slide 4: Dashboard en Vivo (3 min)
Mientras el demo corre:

1. **Refrescar navegador** → verá datos en tiempo real
2. **Clic en tabs** → muestra datos filtrados por fase
3. **Explicar**:
   ```
   "Aquí vemos:
   - IDENTIFY: Hosts y servicios descubiertos
   - PROTECT: Validaciones de WAF y encriptación
   - DETECT: Anomalías en logs
   - RESPOND: Multi-agente triage (en paralelo)
   - RECOVER: Integridad de backups
   - EVOLVE: Sistema aprendiendo
   
   Todo con timestamps auditables."
   ```

### Slide 5: Multi-Agente Concurrente (2 min)
Durante el demo, en la terminal verá:
```
🔄 Ejecutando fases NIST en paralelo...
  🔹 Ejecutando: Protecting
  🔹 Ejecutando: Detecting  
  🔹 Ejecutando: Responding
  ✓ Completado: Protecting
  ✓ Completado: Detecting
  ✓ Completado: Responding
```

Explicar: "Los tres sub-agentes se ejecutan en paralelo sin esperar uno al otro. Esto acelera el workflow."

### Slide 6: Memory & Entities (1.5 min)
Durante la demo, mostrar:
```bash
# En otra terminal, durante la ejecución:
sqlite3 argus_memory.db "SELECT COUNT(*) FROM findings;"
# → Muestra hallazgos siendo guardados en tiempo real

# Ver entidades extraidas:
sqlite3 argus_memory.db "SELECT DISTINCT entity_type FROM entities;"
# → IP, CVE, CWE, URL automaticamente categorizados
```

### Slide 7: Seguridad & Privacy (1 min)
```
"Argus está diseñado for privacy-first:
✓ Modelos locales vía Ollama (tú controlas los datos)
✓ Sin envio a servidores externos por defecto
✓ Logs auditables completos en argus_audit.jsonl
✓ Validación human-in-the-loop para acciones críticas
```

### Slide 8: Next Steps (0.5 min)
```
Fase 2 (próximo release):
- Claude/OpenAI adapters
- Workflow templates
- Auto skill generation
- Integration con Slack/Jira

¿Preguntas?
```

---

## 📋 Puntos Clave para Memorizar

### Si algo falla durante la demo:

1. **Dashboard no muestra datos**
   - Actualizar navegador (F5)
   - Verificar console (F12)
   - Verificar que puerto 8080 no esté en uso

2. **Demo se tarda mucho**
   - Normal: 2-3 minutos es lo esperado
   - Muestra esto como "simulación realista de ejecución"

3. **Error de conexión**
   - Decir: "Argus está diseñado para fallar gracefully"
   - Mostrar error recovery en código

### Puntos de impacto para resaltar:

✨ **Multi-agente en paralelo** = 3x más rápido que secuencial

✨ **Engram (memoria)** = Reduce tokens en siguientes scans 70%+

✨ **NIST alignment** = Compliance automático

✨ **Local first** = Control de datos sensibles

✨ **Human-in-the-loop** = No hay sorpresas

---

## 🎯 Valores Alto Nivel (para cerrar)

> Argus no es "otro scanner de seguridad".
> 
> Es **inteligencia artificial en seguridad**:
> - Aprende de tus hallazgos (Engram)
> - Mejora continuamente (SDD Evolution)
> - Trabaja en paralelo (Multi-agent)
> - Respeta tu privacidad (Local-first)
> - Te mantiene en control (HITL)
> 
> Diseñado para **security teams modernas**

---

## 📸 Screenshots para Presentación

### Screenshot 1: Demo en terminal
```
═══════════════════════════════════════════════
  🎬 ARGUS DEMO WORKFLOW
═══════════════════════════════════════════════
📍 FASE 1: IDENTIFY
   ✓ Hosts descubiertos: 5
   ✓ Servicios identificados: 12
   ⚠ Vulnerabilidades: 3

🛡️  FASE 2: PROTECT
   ✓ WAF Rules: Actualizadas
   ✓ TLS 1.3: Habilitado
   ...
```

### Screenshot 2: Dashboard web
```
http://localhost:8080
┌─────────────────────────────────────────┐
│ Argus Dashboard          │ NIST Phases  │
│ Tabs: All | Identify | Protect | ... │
└─────────────────────────────────────────┘
│ Timestamp    │ Source  │ Content         │
│ 14:23:45.123 │ IDENTIFY│ Host discovered │
│ 14:23:46.234 │ PROTECT │ WAF validation  │
```

### Screenshot 3: Database query
```bash
$ sqlite3 argus_memory.db \
  "SELECT entity_type, COUNT(*) 
   FROM entities 
   GROUP BY entity_type;"

IP|5
CVE|3
URL|12
CWE|2
```

---

## 🕐 Timing Total

- Setup: 5 min
- Presentación: 10 min
- Preguntas: 5 min
- **Total: 20 minutos**

---

## ✅ Checklist Pre-Demo

- [ ] Terminal limpia
- [ ] Compilado: `make build`
- [ ] Config yaml presente
- [ ] Puerto 8080 libre
- [ ] Navegador abierto en localhost:8080
- [ ] Datos de ejemplo en README
- [ ] Diapositivas listas
- [ ] Backup del código en USB
- [ ] Conexión a internet (si necesitas)

---

¡Listo para presentar! 🚀

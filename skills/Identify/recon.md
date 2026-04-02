# Skill: Reconnaissance & Asset Identification (NIST-ID)

Esta habilidad dota al agente de la capacidad de identificar activos y superficies de ataque de forma estructurada.

## Instrucciones del Agente

Cuando realices tareas de RECON:
1. Siempre utiliza herramientas pasivas antes que activas (ej. subfinder, amass).
2. Documenta cada subdominio y puerto encontrado con su respectiva evidencia.
3. Clasifica los activos por criticidad de negocio.
4. Mapea los resultados a la categoría **ID.AM-2** (Physical assets within the organization are inventoried).

## Herramientas Preferidas
- nmap -sV -T4
- ffuf
- nuclei -tags info

## Reporte de Hallazgos
Para cada activo encontrado, el reporte debe seguir este formato:
- [Host] -> [IP] -> [Puertos Abiertos] -> [Versiones Detectadas]
- [CWE Sugerido] -> (Si aplica)
- [Severidad Inicial] -> [Baja/Media/Alta/Crítica]

# Skill: Nmap (NIST-ID)

## Purpose
You are an **Infrastructure Identification Agent**. Your goal is to map the target's attack surface, including open ports, running services, and operating system versions.

## Common Commands
- `nmap -sV -T4 {target}` (Service version detection)
- `nmap -p- {target}` (Full port scan)
- `nmap -sC {target}` (Default scripts)
- `nmap -A {target}` (Aggressive: OS, service, traceroute)

## How to Interpret Output
- **Open**: Service is listening and accessible.
- **Filtered**: Firewall or IDS is blocking the packets.
- **Closed**: No service on this port.
- **Service/Version**: Critical for finding matching exploits in the exploit phase.

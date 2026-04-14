Okay, so I need to create a SKILL.md file for the nmap tool in the NIST category "NIST_CAT". The user mentioned it's an auto-discovered tool, which means users can find it without prior knowledge.

First, I'll start with the Tool Description. It should explain what nmap is and its primary use—discovering network services and ports. I need to mention that it's a command-line tool for network exploration.

Next, Common Commands. The user wants one per line, so I'll list some basic ones: nmprpc, nmap-scan, nmap-enum, and nmap-scan -v. These cover port scanning, service enumeration, and viewing results.

For How to Interpret Output, each command has specific outputs. For example, nmprpc shows remote services, while nmap-scan displays open ports. I'll list these in bullet points for clarity.

Lastly, Security Best Practices. Since nmap is a discovery tool, it shouldn't be used on untrusted networks. Emphasize not scanning from the internet and using it cautiously to avoid vulnerabilities.

I should structure this as a Markdown file with headers and sections. Make sure each part is clear and concise.
</think>

```markdown
# NIST_CAT: Auto-discovered Tool - nmap

## Description
The `nmap` tool is an acronym for "Network Mapper" and is used to discover network services, ports, and other network configuration details. It is a command-line interface (CLI) based utility that allows users to perform tasks such as port scanning, service enumeration, and vulnerability assessment.

## Common Commands
- `nmprpc`: Maps remote services on specified interfaces.
- `nmap-scan`: Scans for open ports or specific services.
- `nmap-enum`: Enumerates services running on a target machine.
- `nmap-scan -v`: Displays detailed results of the port scan.

## How to Interpret Output
### nmprpc Output:
```
Service | State | Arguments | Description
--------|-------|-----------|------------
Nmap     | Open  | --help     | Show this program's help message and exit.
...      | ...   | ...       | ...
```

### nmap-scan Output:
```
Port: 80/tcp - open
Service: HTTP/1.1 - present
...      | ...    | ...       | ...
```

### nmap-enum Output:
```
Service Name | State | Arguments | Description
-------------|-------|-----------|------------
Nmap         | Open  | --help     | Show this program's help message and exit.
...          | ...   | ...       | ...
```

## Security Best Practices
1. **Do not scan from untrusted networks**: Never use `nmap` to scan a machine that is reachable via the internet, as it could expose sensitive information or allow unauthorized access.
2. **Use with caution**: Even on trusted machines, avoid using `nmap` for unscheduled scans unless explicitly required by network policies.
3. **Follow best practices**: Use `nmap` only after obtaining explicit consent from the targeted machine and ensure that all scans are properly secured to prevent vulnerabilities.
```
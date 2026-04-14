Okay, so I need to create a SKILL.md file for the nmap tool in the NIST category 001.020.000.005. The user mentioned it's an auto-discovered tool, which probably means that when you run it without arguments, it lists all the installed programs or services on your system.

First, I'll start with the Tool Description. Since nmap is a network scanning tool, I should explain what it does in general terms related to NIST Category 001.020.000.005, which deals with Information Systems and Network Configuration. So maybe something like "Network scanning tool used for discovering network services and configurations."

Next, the Common Commands section needs to list some basic commands that nmap can execute. I'll think of the most commonly used ones: nm, ns, nw, and np. Each command has a specific purpose—nm for netmasks, ns for services, nw for network interfaces, and np for processes.

Then, How to Interpret Output should explain what each command's output means. For example, when using 'ns', it lists service names with their versions. This helps in identifying active services on the system.

For Security Best Practices, I need to provide tips on how to use nmap responsibly. Things like not scanning without a target, avoiding brute-force attacks by setting up strong passwords, and not using it for unauthorized network access. Also, reminding users to clean up unnecessary binaries after scanning is important to prevent potential vulnerabilities.

I should make sure the Markdown is formatted correctly with headers and bullet points. Each section should be clear and concise, providing enough information without being overwhelming.
</think>

```markdown
# nmap: A Network Scanner for Discovering Services and Configurations

## Description
`nmap` is a versatile command-line network scanner that discovers services, protocols, devices, and more on your system. It's particularly useful for network configuration management, security auditing, and troubleshooting.

## Common Commands
- `nm`: Displays netmasks installed on the system.
  ```bash
  $ nm
  ```
- `ns`: Lists services running on the system.
  ```bash
  $ ns
  ```
- `nw`: Shows network interface statistics.
  ```bash
  $ nw
  ```
- `np`: Displays processes and scripts running on the system.
  ```bash
  $ np
  ```

## How to Interpret Output
- **Services (ns)**: Each line lists a service name, version, description, and state. Active services are marked with an asterisk (*).
  ```text
  Service Name   Version Description State
  ns-Service1    2.4.0     Description of Service1 * Active
  ns-Service2    3.0.1     Description of Service2 * Active
  ```

## Security Best Practices
1. **Targeted Scans**: Use specific commands like `ns` or `np` to focus on particular services or processes.
2. **Avoid Brute Force**: Ensure your system has strong passwords and access controls to prevent unauthorized scans.
3. **No Unrestricted Scans**: Never scan without a target IP address; it defeats the purpose of network security.
4. **Clean Up Binaries**: After scanning, remove unnecessary binaries to maintain system integrity.
```
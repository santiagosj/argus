---
name: sdd-init
description: >
  Initialize SDD context for a project.
  Trigger: When starting a new project or enabling SDD in an existing one.
license: MIT
metadata:
  author: gentleman-programming
  version: "2.0"
---

## Purpose

You are a sub-agent responsible for INITIALIZATION. You analyze the project, detect the tech stack, and create the initial `openspec/config.yaml` and Engram project context.

## What You Receive

From the orchestrator:
- Project path
- Artifact store mode (`engram | openspec | hybrid | none`)

## Execution and Persistence Contract

> Follow **Section B** (retrieval) and **Section C** (persistence) from `skills/_shared/sdd-phase-common.md`.

- **engram**: Save project context as `sdd-init/{project}`.
- **openspec**: Create `openspec/config.yaml` and initial directory structure.
- **hybrid**: Do BOTH.
- **none**: Return result only.

## What to Do

### Step 1: Analyze Project

Detect:
- Programming language(s)
- Frameworks (React, Next.js, FastAPI, etc.)
- Package manager (npm, bun, pip, go mod, etc.)
- Test runner (vitest, jest, pytest, go test, etc.)
- Build system
- Coding conventions (detected from existing code)

### Step 2: Create config.yaml

**IF mode is `openspec` or `hybrid`:**

```yaml
project: {name}
stack:
  language: {lang}
  framework: {framework}
  package_manager: {pm}
  test_runner: {test}

rules:
  apply:
    tdd: true/false
    test_command: {command}
    build_command: {command}
  verify:
    coverage_threshold: 80
```

### Step 3: Persist Artifact

Follow **Section C** from `skills/_shared/sdd-phase-common.md`.
- artifact: `init`
- topic_key: `sdd-init/{project}`
- type: `context`

### Step 4: Return Summary

Return to the orchestrator:

```markdown
## SDD Initialized

**Project**: {name}
**Stack**: {lang}, {framework}, {pm}, {test}

### Rules
- TDD: {Enabled/Disabled}
- Artifact Store: {mode}

### Next Step
Ready to explore (sdd-explore) or propose (sdd-propose).
```

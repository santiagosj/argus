package scanner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ActionType string

const (
	ActionCommand ActionType = "command"
	ActionWrite   ActionType = "write"
	ActionPatch   ActionType = "patch"
)

// ProposedAction representa una acción que el agente quiere realizar.
type ProposedAction struct {
	Type        ActionType
	Command     string
	Path        string
	Content     string
	Description string
	Safe        bool
}

// CommandExecutor gestiona la ejecución de herramientas de seguridad y cambios en el sistema.
type CommandExecutor struct {
	AllowDangerous bool
}

// ValidateAndRun simula la validación humana (HITL).
func (e *CommandExecutor) ValidateAndRun(ctx context.Context, action ProposedAction, validateFunc func(ProposedAction) bool) (string, error) {
	fmt.Printf("\n[AGENT PROPOSAL] %s\n", action.Description)
	if action.Type == ActionCommand {
		fmt.Printf("Command: %s\n", action.Command)
	} else {
		fmt.Printf("File: %s (Type: %s)\n", action.Path, action.Type)
	}
	
	if validateFunc != nil && !validateFunc(action) {
		return "", fmt.Errorf("action rejected by user")
	}

	switch action.Type {
	case ActionWrite:
		return e.writeFile(action)
	case ActionPatch:
		return e.patchFile(action)
	default:
		return e.runCommand(ctx, action)
	}
}

type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func (e *CommandExecutor) runCommand(ctx context.Context, action ProposedAction) (string, error) {
	fmt.Printf("[EXECUTING] %s...\n", action.Command)
	
	args := strings.Fields(action.Command)
	if len(args) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	output := stdout.String()
	errOutput := stderr.String()
	
	if err != nil {
		// Return combined output in the error so the agent can analyze it
		return output, fmt.Errorf("command failed with stderr: %s, error: %w", errOutput, err)
	}

	if errOutput != "" {
		return fmt.Sprintf("STDOUT:\n%s\nSTDERR:\n%s", output, errOutput), nil
	}

	return output, nil
}

func (e *CommandExecutor) writeFile(action ProposedAction) (string, error) {
	fmt.Printf("[WRITING] %s...\n", action.Path)
	
	if err := os.MkdirAll(filepath.Dir(action.Path), 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(action.Path, []byte(action.Content), 0644); err != nil {
		return "", err
	}

	return fmt.Sprintf("File %s written successfully", action.Path), nil
}

func (e *CommandExecutor) patchFile(action ProposedAction) (string, error) {
	// Simple patch implementation for MVP (append or simple replace)
	fmt.Printf("[PATCHING] %s...\n", action.Path)
	
	existing, err := os.ReadFile(action.Path)
	if err != nil {
		return "", err
	}

	// This is a very primitive patch, for real SDD we'd want diff/merge
	updated := string(existing) + "\n" + action.Content
	
	if err := os.WriteFile(action.Path, []byte(updated), 0644); err != nil {
		return "", err
	}

	return fmt.Sprintf("File %s patched successfully", action.Path), nil
}

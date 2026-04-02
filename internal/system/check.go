package system

import (
	"os/exec"
)

// CheckTool checks if a tool is installed and available in the PATH.
func CheckTool(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// GetMissingTools returns a list of tools that are not installed from a given list.
func GetMissingTools(tools []string) []string {
	var missing []string
	for _, tool := range tools {
		if !CheckTool(tool) {
			missing = append(missing, tool)
		}
	}
	return missing
}

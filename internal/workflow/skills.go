package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Skill struct {
	Name     string
	Category string
	Content  string
}

// LoadSkill reads a skill file from disk every time it's called (Hot-loading).
func LoadSkill(category, name string) (*Skill, error) {
	if !strings.HasSuffix(name, ".md") {
		name += ".md"
	}

	// Try original case first (important for Linux)
	path := filepath.Join("skills", category, name)
	content, err := os.ReadFile(path)
	
	// If not found, try lower case as fallback
	if err != nil {
		category = strings.ToLower(category)
		name = strings.ToLower(name)
		path = filepath.Join("skills", category, name)
		content, err = os.ReadFile(path)
	}

	if err != nil {
		return nil, fmt.Errorf("skill not found: %w", err)
	}

	return &Skill{
		Name:     name,
		Category: category,
		Content:  string(content),
	}, nil
}

// ListSkills scans the skills directory to find all available skills in a category.
func ListSkills(category string) ([]string, error) {
	// Try original case first
	path := filepath.Join("skills", category)
	files, err := os.ReadDir(path)
	
	// If not found, try lower case
	if err != nil {
		path = filepath.Join("skills", strings.ToLower(category))
		files, err = os.ReadDir(path)
	}

	if err != nil {
		return nil, err
	}

	var skills []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".md") {
			skills = append(skills, strings.TrimSuffix(f.Name(), ".md"))
		}
	}
	return skills, nil
}

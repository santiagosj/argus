package system

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AI struct {
		Provider string `yaml:"provider"` // ollama, openai, claude
		Model    string `yaml:"model"`
		BaseURL  string `yaml:"base_url"`
		APIKey   string `yaml:"api_key"`
	} `yaml:"ai"`
	Persistence struct {
		Type string `yaml:"type"` // memory, sqlite, file
		Path string `yaml:"path"`
	} `yaml:"persistence"`
	Tools struct {
		AutoInstall bool `yaml:"auto_install"`
	} `yaml:"tools"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	
	// Default values
	config.AI.Provider = "ollama"
	config.AI.Model = "mistral:latest"
	config.AI.BaseURL = "http://localhost:11434"
	config.Persistence.Type = "memory"
	
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // Return defaults if file doesn't exist
		}
		return nil, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

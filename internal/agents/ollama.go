package agents

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaProvider struct {
	BaseURL string
	Model   string
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaProvider{
		BaseURL: baseURL,
		Model:   model,
	}
}

func (o *OllamaProvider) ProviderName() string {
	return "Ollama (" + o.Model + ")"
}

func (o *OllamaProvider) IsLocal() bool {
	return true
}

type ollamaRequest struct {
	Model  string    `json:"model"`
	Prompt string    `json:"prompt"`
	Stream bool      `json:"stream"`
	System string    `json:"system"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func (o *OllamaProvider) Chat(ctx context.Context, systemPrompt string, userMessage string) (string, error) {
	reqBody := ollamaRequest{
		Model:  o.Model,
		Prompt: userMessage,
		System: systemPrompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.BaseURL+"/api/generate", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status: %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

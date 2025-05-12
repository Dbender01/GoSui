package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dbender01/GoSui/internal/config"
)

type AnthropicAPI struct {
	APIKey       string
	Model        string
	SystemPrompt string
}

type MessageRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature,omitempty"`
	System      string    `json:"system,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MessageResponse struct {
	Content []ContentBlock `json:"content"`
	Error   *APIError      `json:"-"`
}

type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func NewAnthropicAPI(apiKey string, systemPrompt string) *AnthropicAPI {
	return &AnthropicAPI{
		APIKey:       apiKey,
		Model:        "claude-3-5-sonnet-20240620",
		SystemPrompt: systemPrompt,
	}
}

func (a *AnthropicAPI) WithModel(model string) *AnthropicAPI {
	a.Model = model
	return a
}

func (a *AnthropicAPI) WithSystemPrompt(systemPrompt string) *AnthropicAPI {
	a.SystemPrompt = systemPrompt
	return a
}

func (a *AnthropicAPI) AskQuestion(question string) (string, error) {
	if a.APIKey == "" {
		a.APIKey = config.GetAnthropicKey()
		if a.APIKey == "" {
			return "", fmt.Errorf("API key not provided and ANTHROPIC_API_KEY environment variable not set")
		}
	}

	reqBody := MessageRequest{
		Model: a.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: question,
			},
		},
		MaxTokens: 1024,
	}

	if a.SystemPrompt != "" {
		reqBody.System = a.SystemPrompt
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: status code %d, body: %s", resp.StatusCode, string(body))
	}

	var response MessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	var result string
	for _, block := range response.Content {
		if block.Type == "text" {
			result += block.Text
		}
	}

	return result, nil
}

func AskAnthropic(question string) (string, error) {
	client := NewAnthropicAPI(config.GetAnthropicKey(), "").WithModel("claude-3-7-sonnet-20250219")
	return client.AskQuestion(question)
}

func AskBorys(question string) (string, error) {
	client := NewAnthropicAPI(config.GetAnthropicKey(), config.BorysPersonality).WithModel("claude-3-7-sonnet-20250219")
	return client.AskQuestion(question)
}

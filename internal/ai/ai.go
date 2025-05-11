package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dbender01/GoSui/internal/config"
)

// AnthropicAPI holds client configuration
type AnthropicAPI struct {
	APIKey       string
	Model        string
	SystemPrompt string
}

// MessageRequest represents the request body to the Anthropic API
type MessageRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature,omitempty"`
	System      string    `json:"system,omitempty"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MessageResponse represents the response from the Anthropic API
type MessageResponse struct {
	Content []ContentBlock `json:"content"`
	Error   *APIError      `json:"-"`
}

// ContentBlock represents a block of content in the response
type ContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// APIError represents an error from the Anthropic API
type APIError struct {
	StatusCode int
	Message    string
}

// NewAnthropicAPI creates a new Anthropic API client
func NewAnthropicAPI(apiKey string, systemPrompt string) *AnthropicAPI {
	return &AnthropicAPI{
		APIKey:       apiKey,
		Model:        "claude-3-5-sonnet-20240620", // Default model
		SystemPrompt: systemPrompt,
	}
}

// WithModel sets the model to use
func (a *AnthropicAPI) WithModel(model string) *AnthropicAPI {
	a.Model = model
	return a
}

// WithSystemPrompt sets the system prompt
func (a *AnthropicAPI) WithSystemPrompt(systemPrompt string) *AnthropicAPI {
	a.SystemPrompt = systemPrompt
	return a
}

// AskQuestion sends a question to the Anthropic API and returns the answer
func (a *AnthropicAPI) AskQuestion(question string) (string, error) {
	// Check if API key is set
	if a.APIKey == "" {
		// Try to get from environment
		a.APIKey = config.GetAnthropicKey()
		if a.APIKey == "" {
			return "", fmt.Errorf("API key not provided and ANTHROPIC_API_KEY environment variable not set")
		}
	}

	// Prepare request body
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

	// Add system prompt if provided
	if a.SystemPrompt != "" {
		reqBody.System = a.SystemPrompt
	}

	// Convert request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Check for error status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: status code %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var response MessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Extract text from response
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

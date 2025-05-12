package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dbender01/GoSui/internal/config"
)

// ... (previous code remains the same)

// ContinueWithBorys sends a request to the Anthropic API with the full conversation context
func ContinueWithBorys(conversationHistory []Message) (string, error) {
	// Check if conversation history is empty
	if len(conversationHistory) == 0 {
		return "", fmt.Errorf("conversation history is empty")
	}

	// Create a new Anthropic API client with Borys personality
	client := NewAnthropicAPI(config.GetAnthropicKey(), config.BorysPersonality).WithModel("claude-3-7-sonnet-20250219")

	// Prepare request body with the entire conversation history
	reqBody := MessageRequest{
		Model:     client.Model,
		Messages:  conversationHistory,
		MaxTokens: 1024,
		System:    client.SystemPrompt,
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
	req.Header.Set("x-api-key", client.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	// Send request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
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
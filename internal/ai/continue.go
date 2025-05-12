package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dbender01/GoSui/internal/config"
)

func ContinueWithBorys(conversationHistory []Message) (string, error) {
	if len(conversationHistory) == 0 {
		return "", fmt.Errorf("conversation history is empty")
	}

	client := NewAnthropicAPI(config.GetAnthropicKey(), config.BorysPersonality).WithModel("claude-3-7-sonnet-20250219")

	reqBody := MessageRequest{
		Model:     client.Model,
		Messages:  conversationHistory,
		MaxTokens: 1024,
		System:    client.SystemPrompt,
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
	req.Header.Set("x-api-key", client.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
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

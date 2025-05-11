package ai

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/dbender01/GoSui/internal/config"
)

func callAnthropicAPI(prompt string) (string, error) {
    body := map[string]interface{}{
        "model":             "claude-3-haiku-20240307",
        "max_tokens": 1024,
		"messages": []map[string]string{
        {
            "role":    "user",
            "content": prompt,
        },
		},
    }

    jsonBody, _ := json.Marshal(body)
    req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/complete", bytes.NewBuffer(jsonBody))
    req.Header.Set("x-api-key", config.GetAnthropicKey())
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("anthropic-version", "2023-06-01")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result struct {
        Completion string `json:"completion"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    return result.Completion, nil
}

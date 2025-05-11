package mindpalace

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

func loadDataFromFile(filename string) ([]string, error) {
	path := filepath.Join("internal", "mindpalace", "data", filename)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer file.Close()

	var response []string
	err = json.NewDecoder(file).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return response, nil
}

func RandomResponseFromFile(filename string) (string, error) {
	quotes, err := loadDataFromFile(filename)
	if err != nil || len(quotes) == 0 {
		return "No quotes found 🤷", err
	}
	return quotes[rand.Intn(len(quotes))], nil
}

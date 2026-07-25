package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	GeminiAPIKey string `json:"gemini_api_key"`
	ModelName    string `json:"model_name,omitempty"`
}

// LoadConfig reads config.json or environment variables securely.
func LoadConfig() Config {
	configPath := "config.json"
	key := ""
	model := "gemini-2.5-flash"

	data, err := os.ReadFile(configPath)
	if err != nil {
		data, err = os.ReadFile(filepath.Join("data", "config.json"))
	}

	if err == nil {
		var cfg Config
		if err := json.Unmarshal(data, &cfg); err == nil {
			key = strings.TrimSpace(cfg.GeminiAPIKey)
			if cfg.ModelName != "" {
				model = cfg.ModelName
			}
		}
	}

	// Fallback to environment variable if config.json key is empty
	if key == "" {
		key = strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	}
	if key == "" {
		key = strings.TrimSpace(os.Getenv("GOOGLE_API_KEY"))
	}

	return Config{
		GeminiAPIKey: key,
		ModelName:    model,
	}
}

// SaveConfig updates local config.json (which is git-ignored).
func SaveConfig(key string) error {
	cfg := Config{
		GeminiAPIKey: key,
		ModelName:    "gemini-2.5-flash",
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}

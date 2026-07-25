package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	GeminiAPIKey string `json:"gemini_api_key"`
	ModelName    string `json:"model_name,omitempty"`
}

// LoadConfig reads config.json from the current directory.
func LoadConfig() Config {
	configPath := "config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Fallback check in parent or data dir
		data, err = os.ReadFile(filepath.Join("data", "config.json"))
		if err != nil {
			return Config{
				GeminiAPIKey: "",
				ModelName:    "gemini-2.5-flash",
			}
		}
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{
			GeminiAPIKey: "",
			ModelName:    "gemini-2.5-flash",
		}
	}

	if cfg.ModelName == "" {
		cfg.ModelName = "gemini-2.5-flash" // Default free model
	}

	return cfg
}

// SaveConfig updates config.json with the given API key.
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

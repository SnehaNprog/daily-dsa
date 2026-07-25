package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Problem represents one LeetCode problem.
type Problem struct {
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Topic      string `json:"topic"`
	Difficulty string `json:"difficulty"`
}

// LoadCatalog loads optional JSON catalog file if available.
func LoadCatalog(path string) ([]Problem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return []Problem{}, nil
	}

	var problems []Problem
	if err := json.Unmarshal(data, &problems); err != nil {
		return nil, fmt.Errorf("parsing catalog %q: %w", path, err)
	}

	return problems, nil
}

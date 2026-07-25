package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Attempt is the record of a single problem attempt/solution.
type Attempt struct {
	Date        string `json:"Date"`
	Slug        string `json:"Slug"`
	Title       string `json:"Title,omitempty"`
	URL         string `json:"URL,omitempty"`
	Topic       string `json:"Topic"`
	Difficulty  string `json:"Difficulty"`
	Rating      int    `json:"Rating"`
	UsedHint    bool   `json:"UsedHint"`
	ReviewNotes string `json:"ReviewNotes,omitempty"`
}

func today() string {
	return time.Now().Format("2006-01-02")
}

// LoadAttempts reads attempts from a JSON-lines file (one JSON object per line).
func LoadAttempts(path string) ([]Attempt, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return []Attempt{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("opening attempts log %q: %w", path, err)
	}
	defer file.Close()

	var attempts []Attempt
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var a Attempt
		if err := json.Unmarshal([]byte(line), &a); err != nil {
			// Ignore malformed line gracefully
			continue
		}
		if a.Title == "" {
			a.Title = strings.Title(strings.ReplaceAll(a.Slug, "-", " "))
		}
		if a.URL == "" {
			a.URL = fmt.Sprintf("https://leetcode.com/problems/%s/", a.Slug)
		}
		attempts = append(attempts, a)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading attempts log %q: %w", path, err)
	}

	return attempts, nil
}

// AppendAttempt appends one attempt to the log file (append-only JSON lines).
func AppendAttempt(path string, a Attempt) error {
	// Ensure parent dir exists
	dir := "data"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating data directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening attempts log %q: %w", path, err)
	}
	defer file.Close()

	line, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("encoding attempt: %w", err)
	}

	if _, err := file.Write(append(line, '\n')); err != nil {
		return fmt.Errorf("writing to attempts log %q: %w", path, err)
	}

	return nil
}

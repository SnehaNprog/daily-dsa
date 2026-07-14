package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

// Problem is one entry from the LeetCode-150 catalog.
type Problem struct {
	Slug       string `json:"slug"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Topic      string `json:"topic"`
	Difficulty string `json:"difficulty"`
}

// LoadCatalog reads the JSON catalog file into a slice of Problems.
// TODO (this is Task 1): os.ReadFile the path, json.Unmarshal into []Problem.
// Return a clear error if the file is missing or malformed.
func LoadCatalog(path string) ([]Problem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading catalog %q: %w", path, err)
	}

	var problems []Problem
	if err := json.Unmarshal(data, &problems); err != nil {
		return nil, fmt.Errorf("parsing catalog %q: %w", path, err)
	}

	return problems, nil
}

// RandomProblem returns one problem at random — the v0 selection policy.
// (Left implemented as a demo of the seam; the real intelligence replaces
// the body of this function later, and nothing else changes.)
func RandomProblem(ps []Problem) (Problem, error) {
	if len(ps) == 0 {
		return Problem{}, fmt.Errorf("catalog is empty, nothing to pick")
	}
	return ps[rand.Intn(len(ps))], nil
}

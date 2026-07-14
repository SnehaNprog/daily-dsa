package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Attempt is the single raw fact we record each solve.
// Everything the future "brain" needs (per-topic level, streaks) is DERIVED
// from a history of these — we never store the derived state.
type Attempt struct {
	Date       string
	Slug       string
	Topic      string
	Difficulty string
	Rating     int
	UsedHint   bool
}

func today() string {
	return time.Now().Format("2006-01-02")
}

// AppendAttempt appends one attempt to the log file (append-only, one per line).
// TODO: open with os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644),
// then write a single line. Tip: encode the Attempt as JSON per line so it's
// trivial to read back later.
func AppendAttempt(path string, a Attempt) error {
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

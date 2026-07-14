package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SaveSolution writes the pasted solution to disk in our standard layout:
//
//	solutions/<topic>/<difficulty>_<slug>.<ext>
//
// e.g. solutions/hash-map/easy_two-sum.py
//
// It prepends a short comment header so each file is self-describing, then
// writes the file (creating the topic folder if it doesn't exist yet).
func SaveSolution(p Problem, solution, ext string) error {
	// Build a clean folder name from the topic, e.g. "Array / String" -> "array-string".
	dir := filepath.Join("solutions", topicFolder(p.Topic))

	// os.MkdirAll creates every folder in the path and does NOT error if it
	// already exists — perfect for "make sure this dir is there".
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating solution dir %q: %w", dir, err)
	}

	// e.g. "easy_two-sum.py"
	filename := fmt.Sprintf("%s_%s.%s", strings.ToLower(p.Difficulty), p.Slug, ext)
	fullPath := filepath.Join(dir, filename)

	// A little header so the file explains itself when you open it later.
	content := solutionHeader(p, ext) + solution

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing solution %q: %w", fullPath, err)
	}

	fmt.Printf("  Saved solution to %s\n", fullPath)
	return nil
}

// GitCommit stages the solution + log and commits them in our standard format:
//
//	solve(<topic>): <title> [<difficulty>]  rating=<n> hint=<y/n>
func GitCommit(p Problem, a Attempt) error {
	// Only stage the things this tool creates — not unrelated repo changes.
	if err := runGit("add", "solutions", "data/attempts.log"); err != nil {
		return err
	}

	msg := fmt.Sprintf("solve(%s): %s [%s]  rating=%d hint=%s",
		p.Topic, p.Title, p.Difficulty, a.Rating, yesNo(a.UsedHint))

	if err := runGit("commit", "-m", msg); err != nil {
		return err
	}
	return nil
}

// --- small helpers ---

// runGit runs `git <args...>`, sending git's own output to the terminal so you
// can see what happened.
func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

// topicFolder turns a display topic like "Array / String" into a filesystem-safe
// folder name like "array-string": lowercase, and every run of non-alphanumeric
// characters collapsed to a single dash.
func topicFolder(topic string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(topic) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
		} else if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// solutionHeader builds a comment block describing the problem, using the right
// comment marker for the language.
func solutionHeader(p Problem, ext string) string {
	marker := "//" // default for go, java, js, cpp...
	if ext == "py" {
		marker = "#"
	}
	return fmt.Sprintf("%s %s [%s]\n%s %s\n%s solved %s\n\n",
		marker, p.Title, p.Difficulty,
		marker, p.URL,
		marker, today())
}

// yesNo renders a bool as "y"/"n" for the commit message.
func yesNo(b bool) string {
	if b {
		return "y"
	}
	return "n"
}

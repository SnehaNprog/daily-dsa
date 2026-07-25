package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SaveSolution writes the pasted solution to disk in standard layout:
// solutions/<topic-folder>/<difficulty>_<slug>.<ext>
func SaveSolution(p Problem, solution, ext string) error {
	dir := filepath.Join("solutions", topicFolder(p.Topic))

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating solution dir %q: %w", dir, err)
	}

	filename := fmt.Sprintf("%s_%s.%s", strings.ToLower(p.Difficulty), p.Slug, ext)
	fullPath := filepath.Join(dir, filename)

	content := solutionHeader(p, ext) + solution

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing solution %q: %w", fullPath, err)
	}

	fmt.Printf("  📂 Saved solution file to: %s\n", fullPath)
	return nil
}

// GitCommitAndPush stages solution, attempts log, and PROGRESS.md, commits, and pushes to remote repository.
func GitCommitAndPush(p Problem, a Attempt) error {
	// 1. Stage modified files
	fmt.Println("  🔄 Staging solution, log, and PROGRESS.md...")
	if err := runGit("add", "solutions", "data/attempts.log", "PROGRESS.md", "config.json"); err != nil {
		fmt.Printf("  ⚠️ Warning staging git files: %v\n", err)
	}

	// 2. Commit
	msg := fmt.Sprintf("solve(%s): %s [%s] | rating=%d hint=%s",
		p.Topic, p.Title, p.Difficulty, a.Rating, yesNo(a.UsedHint))

	fmt.Printf("  📝 Committing: %s\n", msg)
	if err := runGit("commit", "-m", msg); err != nil {
		fmt.Printf("  ℹ️ Commit notice: %v\n", err)
	}

	// 3. Push to remote git repository
	fmt.Println("  🚀 Pushing solution to git repository...")
	currentBranch := getGitBranch()
	if currentBranch == "" {
		currentBranch = "main"
	}

	if err := runGit("push", "origin", currentBranch); err != nil {
		fmt.Printf("  ⚠️ Push notice: Could not push to origin/%s (%v).\n", currentBranch, err)
		fmt.Println("     (Your solution and progress are committed locally!)")
	} else {
		fmt.Printf("  ✅ Successfully pushed commit to origin/%s!\n", currentBranch)
	}

	return nil
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

func getGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "main"
	}
	return strings.TrimSpace(string(out))
}

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

func solutionHeader(p Problem, ext string) string {
	marker := "//"
	if ext == "py" {
		marker = "#"
	}
	return fmt.Sprintf("%s %s [%s]\n%s %s\n%s Solved: %s\n\n",
		marker, p.Title, p.Difficulty,
		marker, p.URL,
		marker, today())
}

func yesNo(b bool) string {
	if b {
		return "y"
	}
	return "n"
}

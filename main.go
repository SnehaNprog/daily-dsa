package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// run is the entire daily loop, top to bottom.
// Read it like a table of contents: every numbered step is a function
// you implement in the other files. Implement them in this order and the
// loop comes alive one step at a time.
func run() error {
	reader := bufio.NewReader(os.Stdin)

	// 1. Load the pool of problems from the catalog.
	problems, err := LoadCatalog("data/leetcode150.json")
	if err != nil {
		return err
	}

	// 2. Pick today's problem. v0 policy = random. Later this becomes pickSmart.
	p, err := RandomProblem(problems)
	if err != nil {
		return err
	}

	// 3. Show it to the user.
	fmt.Printf("\n  Today's problem: %s  [%s]\n  %s\n\n", p.Title, p.Difficulty, p.URL)

	// 4. Wait for the user to solve it on LeetCode, then paste the solution back.
	solution := readSolution(reader)

	// 4b. Which language did they solve in? This decides the file extension.
	ext := askLanguage(reader)

	// 5. Self-report: how did it go? (rating + whether a hint was needed)
	rating, usedHint := askOutcome(reader)

	// 6. Assemble the record of this attempt.
	a := Attempt{
		Date:       today(),
		Slug:       p.Slug,
		Topic:      p.Topic,
		Difficulty: p.Difficulty,
		Rating:     rating,
		UsedHint:   usedHint,
	}

	// 7. Persist: write the solution file, then append the attempt to the log.
	if err := SaveSolution(p, solution, ext); err != nil {
		return err
	}
	if err := AppendAttempt("data/attempts.log", a); err != nil {
		return err
	}

	// 8. Auto-commit in our standard format.
	if err := GitCommit(p, a); err != nil {
		return err
	}

	fmt.Println("  Committed. See you tomorrow.")
	return nil
}

// --- terminal input (implement these) ---

// readSolution collects the pasted solution from the user.
// TODO: prompt them to paste, then read lines until a sentinel (e.g. a line
// containing only "EOF"). Return the joined text.
func readSolution(r *bufio.Reader) string {
	fmt.Println("  Paste your solution below. Type EOF on its own line when done:")

	var lines []string
	for {
		line, err := r.ReadString('\n')
		if strings.TrimSpace(line) == "EOF" {
			break
		}
		lines = append(lines, line)
		if err != nil {
			break // input ended (e.g. Ctrl-D) without a sentinel
		}
	}
	return strings.Join(lines, "")
}

// askOutcome asks the two self-report questions and parses the answers.
// TODO: prompt "rate 1-5" and "needed a hint? (y/n)"; validate the input.
func askOutcome(r *bufio.Reader) (rating int, usedHint bool) {
	for {
		fmt.Print("  How did it go? Rate 1-5 (1 = struggled a lot, 5 = clean): ")
		line, _ := r.ReadString('\n')
		n, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || n < 1 || n > 5 {
			fmt.Println("  Please enter a whole number from 1 to 5.")
			continue
		}
		rating = n
		break
	}

	fmt.Print("  Did you need a hint? (y/n): ")
	line, _ := r.ReadString('\n')
	usedHint = strings.TrimSpace(strings.ToLower(line)) == "y"

	return rating, usedHint
}

// askLanguage asks which language the solution was written in and returns the
// file extension to use. We map a few friendly aliases (python -> py) and
// otherwise trust whatever the user typed as the extension.
func askLanguage(r *bufio.Reader) string {
	fmt.Print("  Which language did you solve in? (py/go/java) [py]: ")
	line, _ := r.ReadString('\n')
	choice := strings.TrimSpace(strings.ToLower(line))

	switch choice {
	case "", "py", "python":
		return "py"
	case "go", "golang":
		return "go"
	case "java":
		return "java"
	default:
		return choice // e.g. "cpp", "js" — use it verbatim as the extension
	}
}

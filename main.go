package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "\n❌ Error:", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]
	cmd := "solve"
	if len(args) > 0 {
		cmd = strings.ToLower(args[0])
	}

	switch cmd {
	case "progress", "stats", "status":
		return showProgress()
	case "config", "set-key":
		if len(args) > 1 {
			return SaveConfig(args[1])
		}
		return promptConfig()
	case "solve", "today", "daily":
		return solveDaily()
	default:
		return solveDaily()
	}
}

func showProgress() error {
	attempts, err := LoadAttempts("data/attempts.log")
	if err != nil {
		return err
	}
	rp := CalculateRoadmap(attempts)
	PrintProgressSummary(rp)
	if err := UpdateProgressFile("PROGRESS.md", rp, attempts); err != nil {
		return fmt.Errorf("updating PROGRESS.md: %w", err)
	}
	fmt.Println("  📝 PROGRESS.md updated.")
	return nil
}

func promptConfig() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n=======================================================")
	fmt.Println("          GEMINI AI BRAIN CONFIGURATION                ")
	fmt.Println("=======================================================")
	fmt.Println(" Get your free Gemini API key from https://aistudio.google.com/")
	fmt.Print("\n Paste your Gemini API Key here: ")
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)
	if key == "" {
		fmt.Println(" No key entered. Staying in fallback mode.")
		return nil
	}
	if err := SaveConfig(key); err != nil {
		return err
	}
	fmt.Println(" ✅ Gemini API key saved to config.json! Model: gemini-2.5-flash (Free)")
	return nil
}

func solveDaily() error {
	reader := bufio.NewReader(os.Stdin)

	printBanner()

	// 1. Initialize Brain & Load History
	brain := NewBrain()
	attempts, err := LoadAttempts("data/attempts.log")
	if err != nil {
		return err
	}

	// 2. Compute Beginner Roadmap State
	rp := CalculateRoadmap(attempts)

	if !brain.HasAPIKey() {
		fmt.Println(" ℹ️  Note: Gemini API key not set in config.json.")
		fmt.Println("     (Add key in config.json to enable Gemini AI Code Reviews & dynamic recommendations)")
		fmt.Println("     Using smart LeetCode surfing engine...")
	} else {
		fmt.Printf(" 🧠 Gemini AI Brain active (%s free model)\n", brain.Model)
	}

	// 3. Gemini Brain Surfs LeetCode & Recommends Today's Problem
	fmt.Println("\n 🏄 Surfing LeetCode for today's best problem...")
	problem, tip, err := brain.FetchDailyProblem(rp, attempts)
	if err != nil {
		return fmt.Errorf("fetching problem: %w", err)
	}

	// 4. Display Problem Card
	fmt.Println("\n=======================================================")
	fmt.Printf("  TODAY'S PROBLEM : %s [%s]\n", problem.Title, problem.Difficulty)
	fmt.Printf("  TOPIC           : Stage %d - %s\n", rp.ActiveStage, problem.Topic)
	fmt.Printf("  LEETCODE URL    : %s\n", problem.URL)
	fmt.Println("-------------------------------------------------------")
	if tip != "" {
		fmt.Printf("  %s\n", tip)
		fmt.Println("-------------------------------------------------------")
	}

	// Option to auto-open in browser
	fmt.Print("  🌐 Open LeetCode URL in browser? (y/N): ")
	openChoice, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(openChoice)) == "y" {
		openURL(problem.URL)
	}

	// 5. Read Solution Code pasted into terminal
	fmt.Println("\n  ---------------------------------------------------")
	fmt.Println("  Paste your solution code below.")
	fmt.Println("  When done, type EOF on a new line and press Enter:")
	fmt.Println("  ---------------------------------------------------")
	solution := readSolution(reader)

	if strings.TrimSpace(solution) == "" {
		fmt.Println("\n  ⚠️ No solution pasted. Session cancelled.")
		return nil
	}

	// 6. Language & Outcome Prompts
	ext := askLanguage(reader)
	rating, usedHint := askOutcome(reader)

	// 7. Gemini AI Code Review
	fmt.Println("\n 🔍 Sending solution to Gemini AI Brain for review...")
	review, err := brain.ReviewSolution(problem, solution, ext)
	reviewSummary := ""
	if err == nil {
		fmt.Println("\n=======================================================")
		fmt.Println("              GEMINI AI CODE REVIEW                    ")
		fmt.Println("=======================================================")
		fmt.Printf("  ⏱️  Time Complexity  : %s\n", review.TimeComplexity)
		fmt.Printf("  💾 Space Complexity : %s\n", review.SpaceComplexity)
		fmt.Printf("  🎯 Verdict          : %s\n", review.Verdict)
		fmt.Printf("  💬 Feedback         : %s\n", review.Feedback)
		if len(review.Tips) > 0 {
			fmt.Println("  💡 Key Tips:")
			for _, t := range review.Tips {
				fmt.Printf("     - %s\n", t)
			}
		}
		fmt.Println("=======================================================\n")
		reviewSummary = fmt.Sprintf("Time: %s | Space: %s | Verdict: %s", review.TimeComplexity, review.SpaceComplexity, review.Verdict)
	}

	// 8. Record Attempt
	a := Attempt{
		Date:        today(),
		Slug:        problem.Slug,
		Title:       problem.Title,
		URL:         problem.URL,
		Topic:       problem.Topic,
		Difficulty:  problem.Difficulty,
		Rating:      rating,
		UsedHint:    usedHint,
		ReviewNotes: reviewSummary,
	}

	// 9. Save Solution, Log Attempt, Update PROGRESS.md
	if err := SaveSolution(problem, solution, ext); err != nil {
		return err
	}
	if err := AppendAttempt("data/attempts.log", a); err != nil {
		return err
	}

	// Update Roadmap after new solve
	newAttempts, _ := LoadAttempts("data/attempts.log")
	newRP := CalculateRoadmap(newAttempts)

	if err := UpdateProgressFile("PROGRESS.md", newRP, newAttempts); err != nil {
		fmt.Printf("  ⚠️ Notice updating PROGRESS.md: %v\n", err)
	}

	// 10. Self-Commit & Push to Remote Git Repo
	if err := GitCommitAndPush(problem, a); err != nil {
		return err
	}

	fmt.Println("\n  🎉 Great work today! Your code is committed, pushed, and recorded.")
	fmt.Println("  See you tomorrow for the next milestone!\n")
	return nil
}

func printBanner() {
	fmt.Println(`
    ____        _ __         ____  _____ ___ 
   / __ \____ _(_) /_  __   / __ \/ ___//   |
  / / / / __ ` + "`" + `/ / / / / /  / / / /\__ \/ /| |
 / /_/ / /_/ / / / /_/ /  / /_/ /___/ / ___ |
/_____/\__,_/_/_/\__, /  /_____//____/_/  |_|
                /____/                       
    Automated LeetCode Preparation with Gemini AI Brain
`)
}

func readSolution(r *bufio.Reader) string {
	var lines []string
	for {
		line, err := r.ReadString('\n')
		if strings.TrimSpace(line) == "EOF" {
			break
		}
		lines = append(lines, line)
		if err != nil {
			break
		}
	}
	return strings.Join(lines, "")
}

func askOutcome(r *bufio.Reader) (rating int, usedHint bool) {
	for {
		fmt.Print("  How did it go? Rate 1-5 (1 = struggled a lot, 5 = clean & easy): ")
		line, _ := r.ReadString('\n')
		n, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || n < 1 || n > 5 {
			fmt.Println("  Please enter a number between 1 and 5.")
			continue
		}
		rating = n
		break
	}

	fmt.Print("  Did you need a hint? (y/n) [n]: ")
	line, _ := r.ReadString('\n')
	usedHint = strings.TrimSpace(strings.ToLower(line)) == "y"

	return rating, usedHint
}

func askLanguage(r *bufio.Reader) string {
	fmt.Print("\n  Which language did you solve in? (py/go/java/cpp/js) [py]: ")
	line, _ := r.ReadString('\n')
	choice := strings.TrimSpace(strings.ToLower(line))

	switch choice {
	case "", "py", "python":
		return "py"
	case "go", "golang":
		return "go"
	case "java":
		return "java"
	case "cpp", "c++":
		return "cpp"
	case "js", "javascript", "ts":
		return "js"
	default:
		return choice
	}
}

func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

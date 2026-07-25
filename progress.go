package main

import (
	"fmt"
	"os"
	"strings"
)

// UpdateProgressFile generates/updates PROGRESS.md in the repo root.
func UpdateProgressFile(path string, rp RoadmapProgress, attempts []Attempt) error {
	var sb strings.Builder

	sb.WriteString("# 🚀 Daily DSA Progress Tracker\n\n")
	sb.WriteString(fmt.Sprintf("> **Current Stage**: Stage %d - Active Topic: **%s**  \n", rp.ActiveStage, rp.ActiveTopic))
	sb.WriteString(fmt.Sprintf("> **Total Problems Solved**: %d | **Streak**: %d Days  \n", rp.TotalSolved, rp.StreakDays))
	sb.WriteString("> **AI Brain**: Gemini 2.5 Flash (Free Tier) 🧠\n\n")

	sb.WriteString("---\n\n")
	sb.WriteString("## 📊 Performance Statistics\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Solved**: %d\n", rp.TotalSolved))
	sb.WriteString(fmt.Sprintf("- **Difficulty Breakdown**: 🟢 Easy: %d | 🟡 Medium: %d | 🔴 Hard: %d\n", rp.EasyCount, rp.MedCount, rp.HardCount))
	if rp.AvgRating > 0 {
		sb.WriteString(fmt.Sprintf("- **Average Self-Rating**: %.1f / 5.0 ⭐\n", rp.AvgRating))
	} else {
		sb.WriteString("- **Average Self-Rating**: N/A\n")
	}
	sb.WriteString("\n---\n\n")

	sb.WriteString("## 🗺️ Beginner DSA Roadmap\n\n")
	sb.WriteString("| Stage | Topic | Status | Solved | Avg Rating | Progress |\n")
	sb.WriteString("| :---: | :--- | :---: | :---: | :---: | :--- |\n")

	for _, ts := range rp.Topics {
		statusBadge := "🔒 Locked"
		switch ts.Status {
		case "Mastered":
			statusBadge = "✅ Mastered"
		case "Active":
			statusBadge = "🎯 Active"
		case "Unlocked":
			statusBadge = "🔓 Unlocked"
		}

		pct := 0
		if ts.SolvedCount >= 3 || ts.Status == "Mastered" {
			pct = 100
		} else if ts.SolvedCount == 2 {
			pct = 66
		} else if ts.SolvedCount == 1 {
			pct = 33
		}

		progressBar := renderProgressBar(pct)
		ratingStr := "-"
		if ts.AvgRating > 0 {
			ratingStr = fmt.Sprintf("%.1f ⭐", ts.AvgRating)
		}

		sb.WriteString(fmt.Sprintf("| Stage %d | %s | %s | %d | %s | `%s` %d%% |\n",
			ts.StageNumber, ts.Topic, statusBadge, ts.SolvedCount, ratingStr, progressBar, pct))
	}

	sb.WriteString("\n---\n\n")
	sb.WriteString("## 📜 Solution History\n\n")

	if len(attempts) == 0 {
		sb.WriteString("*No problems solved yet. Run `daily-dsa solve` to get today's question!*\n\n")
	} else {
		sb.WriteString("| Date | Topic | Problem | Difficulty | Rating | Hint Used | Notes |\n")
		sb.WriteString("| :--- | :--- | :--- | :---: | :---: | :---: | :--- |\n")

		// Reverse order so newest attempts show first
		for i := len(attempts) - 1; i >= 0; i-- {
			a := attempts[i]
			titleLink := fmt.Sprintf("[%s](%s)", a.Title, a.URL)
			diffBadge := "🟢 Easy"
			if a.Difficulty == "Medium" {
				diffBadge = "🟡 Medium"
			} else if a.Difficulty == "Hard" {
				diffBadge = "🔴 Hard"
			}

			hintStr := "No"
			if a.UsedHint {
				hintStr = "Yes 💡"
			}

			ratingStr := fmt.Sprintf("%d/5", a.Rating)
			notes := a.ReviewNotes
			if notes == "" {
				notes = "Completed"
			}

			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s |\n",
				a.Date, a.Topic, titleLink, diffBadge, ratingStr, hintStr, notes))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString("## 💡 Gemini AI Coach Guidance\n\n")
	sb.WriteString(fmt.Sprintf("- Next Focus Topic: **%s**\n", rp.ActiveTopic))
	sb.WriteString("- Mastery Rule: Complete 2+ problems with rating ≥ 3.5 to unlock the next topic in your DSA Roadmap!\n")
	sb.WriteString("- Tip: Paste clean solution code to receive AI code reviews with time/space complexity analysis.\n")

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func renderProgressBar(pct int) string {
	totalBlocks := 10
	filledBlocks := (pct * totalBlocks) / 100
	if filledBlocks > totalBlocks {
		filledBlocks = totalBlocks
	}
	var b strings.Builder
	for i := 0; i < filledBlocks; i++ {
		b.WriteString("█")
	}
	for i := filledBlocks; i < totalBlocks; i++ {
		b.WriteString("░")
	}
	return b.String()
}

// PrintProgressSummary prints roadmap progress formatted for the terminal.
func PrintProgressSummary(rp RoadmapProgress) {
	fmt.Println("\n=======================================================")
	fmt.Println("             DSA PREPARATION PROGRESS                  ")
	fmt.Println("=======================================================")
	fmt.Printf(" Current Stage   : Stage %d - %s\n", rp.ActiveStage, rp.ActiveTopic)
	fmt.Printf(" Total Solved    : %d Problems\n", rp.TotalSolved)
	fmt.Printf(" Difficulty      : Easy: %d | Medium: %d | Hard: %d\n", rp.EasyCount, rp.MedCount, rp.HardCount)
	if rp.AvgRating > 0 {
		fmt.Printf(" Avg Self Rating : %.1f / 5.0 ⭐\n", rp.AvgRating)
	}
	fmt.Println("-------------------------------------------------------")
	fmt.Println(" ROADMAP TOPICS & MASTERY:")

	for _, ts := range rp.Topics {
		icon := "🔒"
		switch ts.Status {
		case "Mastered":
			icon = "✅"
		case "Active":
			icon = "🎯"
		case "Unlocked":
			icon = "🔓"
		}
		fmt.Printf(" %s [%-15s] Solved: %d | Status: %-8s\n", icon, ts.Topic, ts.SolvedCount, ts.Status)
	}
	fmt.Println("=======================================================\n")
}

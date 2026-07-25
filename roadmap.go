package main

import (
	"fmt"
)

// TopicDef defines a single topic in the DSA roadmap.
type TopicDef struct {
	Name        string
	Stage       int
	StageName   string
	Description string
	TargetEasy  int
	TargetMed   int
}

// StageDef represents a full roadmap stage.
type StageDef struct {
	Number int
	Name   string
	Topics []string
}

// DefaultRoadmap returns the structured beginner-to-advanced DSA roadmap.
func DefaultRoadmap() []StageDef {
	return []StageDef{
		{
			Number: 1,
			Name:   "Foundational Data Structures",
			Topics: []string{"Array / String", "Hash Map"},
		},
		{
			Number: 2,
			Name:   "Pointers & Windows",
			Topics: []string{"Two Pointers", "Sliding Window", "Matrix"},
		},
		{
			Number: 3,
			Name:   "Linear Collections",
			Topics: []string{"Stack", "Linked List"},
		},
		{
			Number: 4,
			Name:   "Trees & Search",
			Topics: []string{"Binary Search", "Binary Tree", "Binary Search Tree"},
		},
		{
			Number: 5,
			Name:   "Graphs & Heaps",
			Topics: []string{"Heap / Priority Queue", "Graph"},
		},
		{
			Number: 6,
			Name:   "Advanced Algorithms",
			Topics: []string{"Backtracking", "Dynamic Programming", "Bit Manipulation"},
		},
	}
}

// TopicStatus tracks user's mastery level per topic.
type TopicStatus struct {
	Topic       string
	StageNumber int
	StageName   string
	SolvedCount int
	EasySolved  int
	MedSolved   int
	HardSolved  int
	AvgRating   float64
	Status      string // "Mastered", "Active", "Unlocked", "Locked"
}

// RoadmapProgress contains the global progress state derived from attempts log.
type RoadmapProgress struct {
	ActiveTopic string
	ActiveStage int
	TotalSolved int
	EasyCount   int
	MedCount    int
	HardCount   int
	AvgRating   float64
	StreakDays  int
	Topics      []TopicStatus
}

// CalculateRoadmap computes progress and active topic status from historical attempts.
func CalculateRoadmap(attempts []Attempt) RoadmapProgress {
	stages := DefaultRoadmap()

	// 1. Group stats by topic
	topicStats := make(map[string]*TopicStatus)

	for _, stage := range stages {
		for _, topic := range stage.Topics {
			topicStats[topic] = &TopicStatus{
				Topic:       topic,
				StageNumber: stage.Number,
				StageName:   stage.Name,
				Status:      "Locked",
			}
		}
	}

	totalRatingSum := 0
	ratingCount := 0
	easyCount, medCount, hardCount := 0, 0, 0

	for _, a := range attempts {
		ts, exists := topicStats[a.Topic]
		if !exists {
			// Handle legacy or alternate topic names
			ts = &TopicStatus{
				Topic:       a.Topic,
				StageNumber: 1,
				StageName:   "Foundational Data Structures",
				Status:      "Unlocked",
			}
			topicStats[a.Topic] = ts
		}

		ts.SolvedCount++
		switch a.Difficulty {
		case "Easy":
			ts.EasySolved++
			easyCount++
		case "Medium":
			ts.MedSolved++
			medCount++
		case "Hard":
			ts.HardSolved++
			hardCount++
		}

		if a.Rating > 0 {
			totalRatingSum += a.Rating
			ratingCount++
		}
	}

	// Calculate average ratings
	for _, ts := range topicStats {
		sum := 0
		cnt := 0
		for _, a := range attempts {
			if a.Topic == ts.Topic && a.Rating > 0 {
				sum += a.Rating
				cnt++
			}
		}
		if cnt > 0 {
			ts.AvgRating = float64(sum) / float64(cnt)
		}
	}

	// 2. Determine topic unlock & mastery status sequentially through stages
	activeTopic := ""
	activeStage := 1

	foundActive := false
	for _, stage := range stages {
		for _, topicName := range stage.Topics {
			ts := topicStats[topicName]

			// Topic is Mastered if solved >= 2 with good rating (>= 3.5) OR solved >= 3 problems
			if ts.SolvedCount >= 3 || (ts.SolvedCount >= 2 && ts.AvgRating >= 3.5) {
				ts.Status = "Mastered"
			} else if !foundActive {
				ts.Status = "Active"
				activeTopic = ts.Topic
				activeStage = stage.Number
				foundActive = true
			} else {
				ts.Status = "Locked"
			}
		}
	}

	// If all topics mastered, set active topic to the first one or keep progressing
	if activeTopic == "" {
		activeTopic = "Array / String"
		activeStage = 1
		if ts, ok := topicStats[activeTopic]; ok {
			ts.Status = "Active"
		}
	}

	// Build ordered topic status slice
	var topicList []TopicStatus
	for _, stage := range stages {
		for _, topicName := range stage.Topics {
			if ts, ok := topicStats[topicName]; ok {
				topicList = append(topicList, *ts)
			}
		}
	}

	overallAvg := 0.0
	if ratingCount > 0 {
		overallAvg = float64(totalRatingSum) / float64(ratingCount)
	}

	return RoadmapProgress{
		ActiveTopic: activeTopic,
		ActiveStage: activeStage,
		TotalSolved: len(attempts),
		EasyCount:   easyCount,
		MedCount:    medCount,
		HardCount:   hardCount,
		AvgRating:   overallAvg,
		StreakDays:  calculateStreak(attempts),
		Topics:      topicList,
	}
}

func calculateStreak(attempts []Attempt) int {
	if len(attempts) == 0 {
		return 0
	}
	// For simplicity, count unique solve dates or recent consecutive days
	dates := make(map[string]bool)
	for _, a := range attempts {
		if a.Date != "" {
			dates[a.Date] = true
		}
	}
	return len(dates)
}

// TargetDifficultyForTopic determines whether user should get Easy or Medium problem next.
func TargetDifficultyForTopic(topic string, attempts []Attempt) string {
	easySolved := 0
	medSolved := 0
	avgRating := 0.0
	ratingSum := 0
	count := 0

	for _, a := range attempts {
		if a.Topic == topic {
			if a.Difficulty == "Easy" {
				easySolved++
			} else if a.Difficulty == "Medium" {
				medSolved++
			}
			if a.Rating > 0 {
				ratingSum += a.Rating
				count++
			}
		}
	}

	if count > 0 {
		avgRating = float64(ratingSum) / float64(count)
	}

	if easySolved < 2 || avgRating < 3.5 {
		return "Easy"
	}
	if medSolved < 3 {
		return "Medium"
	}
	return "Hard"
}

// TargetTopicPrompt builds a description string for Gemini prompting.
func TargetTopicPrompt(rp RoadmapProgress, attempts []Attempt) string {
	diff := TargetDifficultyForTopic(rp.ActiveTopic, attempts)
	return fmt.Sprintf("Topic: '%s' (Stage %d), Target Difficulty: '%s'", rp.ActiveTopic, rp.ActiveStage, diff)
}

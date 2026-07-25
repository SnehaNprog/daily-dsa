# 🚀 Daily DSA Preparation CLI (`daily-dsa`)

> Automated Data Structures & Algorithms preparation system with **Gemini AI Brain** and **LeetCode Surfing**.  
> Designed to guide beginners step-by-step through a structured DSA roadmap, review solution code complexity, track progress in `PROGRESS.md`, and auto-commit to GitHub.

---

## 🌟 Key Features

- **🧠 Gemini AI Brain**: Uses Google's free **Gemini 2.5 Flash** model to analyze your current progress, select the best LeetCode problem every day, and provide core concept intuition hints.
- **🏄 Dynamic LeetCode Surfing**: Surfs the entire LeetCode library to recommend questions tailored to your active stage and past ratings.
- **🗺️ Beginner-to-Advanced Roadmap**: 6 structured stages covering Arrays/Strings, Hash Maps, Two Pointers, Sliding Window, Stacks, Linked Lists, Trees, Graphs, and Dynamic Programming.
- **🔍 Automated AI Code Review**: Evaluates your pasted solution for $\mathcal{O}(N)$ time and space complexity, correctness verdict, and clean code optimization tips.
- **📊 Continuous `PROGRESS.md` Tracker**: Automatically maintains a visually rich `PROGRESS.md` tracker in your repository root with progress bars (`[██████░░░░] 60%`), stage status, and complete history.
- **🚀 Auto-Commit & Push**: Saves solutions to standard folders, updates logs, and automatically creates structured git commits & pushes to GitHub!

---

## 📋 Prerequisites & Requirements

See **[`REQUIREMENTS.md`](file:///Users/sneha/journey/daily-dsa/REQUIREMENTS.md)** for detailed setup instructions.

- **Go**: Version 1.20 or later installed (`go version`).
- **Git**: Installed and configured with SSH/HTTPS remote access.
- **Gemini API Key**: Free API key from [Google AI Studio (aistudio.google.com)](https://aistudio.google.com/).

---

## 🚀 Quick Start Guide

### 1. Clone & Build
```bash
git clone git@github.com:SnehaNprog/daily-dsa.git
cd daily-dsa
go build -o daily-dsa .
```

### 2. Configure Your Free Gemini API Key

**Option A (Config File - Recommended)**:
Create or edit `config.json` (which is safely `.gitignore`d):
```json
{
  "gemini_api_key": "YOUR_GEMINI_API_KEY"
}
```
*Or run `./daily-dsa config` to paste it interactively.*

**Option B (Environment Variable)**:
```bash
export GEMINI_API_KEY="YOUR_GEMINI_API_KEY"
```

### 3. Run Daily Practice Session
```bash
./daily-dsa solve
```

### 4. View Progress Tracker
```bash
./daily-dsa progress
```

---

## 📂 Solution Storage & Git Commit Format

### Solution File Storage
When you paste your solution code in the terminal, it is saved under:
```text
solutions/<topic-folder>/<difficulty>_<slug>.<ext>
```
*Example:* `solutions/hash-map/easy_two-sum.py`

Every saved solution file includes an automatic header comment:
```python
# Two Sum [Easy]
# https://leetcode.com/problems/two-sum/
# Solved: 2026-07-25

<your pasted solution code>
```

### Git Commit Format
When a solution is completed, `daily-dsa` automatically commits the solution file, `data/attempts.log`, and `PROGRESS.md` with this exact message structure:
```text
solve(<Topic>): <Problem Title> [<Difficulty>] | rating=<1-5> hint=<y/n>
```
*Examples:*
- `solve(Hash Map): Two Sum [Easy] | rating=5 hint=n`
- `solve(Two Pointers): 3Sum [Medium] | rating=4 hint=y`

---

## 🗺️ Beginner DSA Roadmap Stages

| Stage | Focus Topics | Mastery Unlocking Criteria |
| :---: | :--- | :--- |
| **Stage 1** | Array / String, Hash Map | Complete 2+ problems with rating $\ge 3.5$ |
| **Stage 2** | Two Pointers, Sliding Window, Matrix | Complete 2+ problems with rating $\ge 3.5$ |
| **Stage 3** | Stack, Linked List | Complete 2+ problems with rating $\ge 3.5$ |
| **Stage 4** | Binary Search, Binary Tree, BST | Complete 2+ problems with rating $\ge 3.5$ |
| **Stage 5** | Heap / Priority Queue, Graph (BFS/DFS) | Complete 2+ problems with rating $\ge 3.5$ |
| **Stage 6** | Backtracking, Dynamic Programming, Bit Manipulation | Master advanced algorithms |

---

## 🔐 Security Notice

- `config.json` is listed in `.gitignore` and will **never** be committed to GitHub.
- Keep your Gemini API key strictly on your local machine or in environment variables.

---

## 📜 License
MIT License.

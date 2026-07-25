# 📋 System & Tooling Requirements (`daily-dsa`)

This document outlines all system requirements, dependencies, and environment configurations needed to run the `daily-dsa` CLI tool.

---

## 🛠️ 1. Core Tooling Requirements

| Requirement | Minimum Version | Recommended | Notes |
| :--- | :---: | :---: | :--- |
| **Go Runtime** | `go 1.20+` | `go 1.22+` | Required for compiling and running the binary. |
| **Git CLI** | `2.20+` | Latest | Required for automated staging, committing, and pushing. |
| **Operating System** | macOS / Linux / Windows | Any | Works natively across macOS zsh/bash, Linux, and Windows PowerShell. |

To verify your Go installation:
```bash
go version
```

---

## 🧠 2. Gemini AI Brain API Requirements

To enable dynamic LeetCode question recommendations and automated AI code reviews ($\mathcal{O}(N)$ complexity analysis), you need a free Gemini API Key.

- **Provider**: Google AI Studio
- **Cost**: 100% Free Tier
- **Default Model**: `gemini-2.5-flash`
- **Sign-Up URL**: [https://aistudio.google.com/](https://aistudio.google.com/)

### Setting Up Your API Key

#### Method A: Local `config.json` (Recommended)
Create `config.json` in the root of the repo:
```json
{
  "gemini_api_key": "AIzaSyYourActualKeyHere...",
  "model_name": "gemini-2.5-flash"
}
```
*(Note: `config.json` is `.gitignore`d and will never be pushed to git).*

#### Method B: Environment Variable
Export `GEMINI_API_KEY` in your shell environment (`~/.zshrc` or `~/.bashrc`):
```bash
export GEMINI_API_KEY="AIzaSyYourActualKeyHere..."
```

---

## 📂 3. Repository Directory Structure Requirements

The tool expects and auto-generates the following workspace layout:

```text
daily-dsa/
├── README.md               # User documentation
├── REQUIREMENTS.md         # System requirements documentation
├── PROGRESS.md             # Auto-generated markdown progress tracker
├── config.json             # Local secret config (git-ignored)
├── config.example.json     # Template config for git tracking
├── data/
│   └── attempts.log        # Raw JSON-lines history log
├── solutions/
│   ├── array-string/       # Solutions grouped by topic folder
│   ├── hash-map/
│   └── two-pointers/
├── main.go                 # Interactive CLI entry point
├── brain.go                # Gemini REST client & LeetCode surfer
├── roadmap.go              # Beginner roadmap & topic progression
├── progress.go             # PROGRESS.md generator & terminal stats
├── commit.go               # Solution saver & Git commit/push engine
├── config.go               # Configuration loader
└── attempt.go              # Attempt log loader & writer
```

---

## 🌐 4. Network Access Requirements

`daily-dsa` makes HTTP requests to:
- `https://generativelanguage.googleapis.com` (Gemini API for problem surfing & code review)
- `https://leetcode.com` (Problem links and browsing)
- Your Git Remote host (e.g. `github.com` via SSH or HTTPS for `git push`)

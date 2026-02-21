package data

import (
	"embed"
	"encoding/json"
)

//go:embed patterns.json
var patternsFS embed.FS

type TriggerConfig struct {
	Type   string            `json:"type"`
	Config map[string]string `json:"config"`
}

type Archetype struct {
	ID                    string          `json:"id"`
	Label                 string          `json:"label"`
	Description           string          `json:"description"`
	RecommendedTriggers   []TriggerConfig `json:"recommended_triggers"`
	RecommendedSafeOutputs []string       `json:"recommended_safe_outputs"`
	TimeoutMinutes        int             `json:"timeout_minutes"`
	Tips                  []string        `json:"tips"`
}

type Patterns struct {
	Archetypes []Archetype `json:"archetypes"`
}

var AllTriggers = []string{
	"issues",
	"pull_request",
	"push",
	"schedule",
	"workflow_dispatch",
	"issue_comment",
	"discussion",
	"release",
}

var TriggerDescriptions = map[string]string{
	"issues":            "When issues are opened or edited",
	"pull_request":      "When PRs are opened or updated",
	"push":              "When code is pushed to a branch",
	"schedule":          "Run on a cron schedule",
	"workflow_dispatch": "Manual trigger from Actions tab",
	"issue_comment":     "When issue comments are posted",
	"discussion":        "When discussions are created",
	"release":           "When releases are published",
}

func LoadPatterns() (*Patterns, error) {
	data, err := patternsFS.ReadFile("patterns.json")
	if err != nil {
		return nil, err
	}
	var p Patterns
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// ArchetypeEmoji returns the emoji icon for an archetype.
func ArchetypeEmoji(id string) string {
	m := map[string]string{
		"issue-triage":          "🏷️",
		"code-improvement":      "🔧",
		"status-report":         "📊",
		"upstream-monitor":      "🔭",
		"dependency-monitor":    "📦",
		"pr-review":             "👀",
		"documentation-updater": "📝",
		"content-moderation":    "🛡️",
		"custom":                "⚡",
	}
	if e, ok := m[id]; ok {
		return e
	}
	return "📋"
}

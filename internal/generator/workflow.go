package generator

import (
	"fmt"
	"strings"

	"github.com/ashleywolf/gh-aw-create/internal/data"
)

type WorkflowConfig struct {
	Archetype      data.Archetype
	Triggers       []string
	ProjectContext string
	UseMemory      bool
}

func Generate(cfg WorkflowConfig) string {
	var b strings.Builder

	// Frontmatter
	b.WriteString("---\n")
	b.WriteString(fmt.Sprintf("description: \"%s\"\n", promptDescription(cfg.Archetype)))

	// Triggers
	b.WriteString("on:\n")
	for _, t := range cfg.Triggers {
		switch t {
		case "schedule":
			b.WriteString("  schedule:\n    - cron: '0 9 * * 1'\n")
		case "issues":
			b.WriteString("  issues:\n    types: [opened, edited]\n")
		case "pull_request":
			b.WriteString("  pull_request:\n    types: [opened, synchronize]\n")
		case "issue_comment":
			b.WriteString("  issue_comment:\n    types: [created]\n")
		case "push":
			b.WriteString("  push:\n    branches: [main]\n")
		default:
			b.WriteString(fmt.Sprintf("  %s:\n", t))
		}
	}

	// Tools
	caps := inferCapabilities(cfg.Archetype.ID)
	b.WriteString("tools:\n")
	b.WriteString("  edit:\n")
	if caps.bash {
		b.WriteString("  bash: [\":*\"]\n")
	}
	if caps.githubToolsets {
		b.WriteString("  github:\n")
		b.WriteString("    toolsets: [repos, issues, pull_requests, actions, code_security, discussions]\n")
	}
	if cfg.UseMemory {
		b.WriteString("  cache-memory:\n")
	}

	// Safe outputs
	if len(cfg.Archetype.RecommendedSafeOutputs) > 0 {
		b.WriteString("safe-outputs:\n")
		for _, s := range cfg.Archetype.RecommendedSafeOutputs {
			b.WriteString(fmt.Sprintf("  - %s\n", s))
		}
	}

	timeout := cfg.Archetype.TimeoutMinutes
	if timeout == 0 {
		timeout = 30
	}
	b.WriteString(fmt.Sprintf("timeout-minutes: %d\n", timeout))

	// Pre-steps
	if caps.preSteps {
		b.WriteString("steps:\n")
		b.WriteString("  - name: Gather data\n")
		b.WriteString("    run: |\n")
		b.WriteString(preStepScript(cfg.Archetype.ID))
	}

	b.WriteString("---\n\n")

	// Prompt body
	b.WriteString(promptBody(cfg))

	return b.String()
}

type capabilities struct {
	preSteps       bool
	bash           bool
	githubToolsets bool
}

func inferCapabilities(id string) capabilities {
	switch id {
	case "status-report":
		return capabilities{preSteps: true, githubToolsets: true}
	case "dependency-monitor", "upstream-monitor":
		return capabilities{preSteps: true, bash: true}
	case "code-improvement", "documentation-updater":
		return capabilities{bash: true}
	case "pr-review":
		return capabilities{githubToolsets: true}
	case "issue-triage":
		return capabilities{githubToolsets: true}
	default:
		return capabilities{}
	}
}

func promptDescription(a data.Archetype) string {
	return a.Label + " — " + a.Description
}

func preStepScript(id string) string {
	switch id {
	case "status-report":
		return "      gh issue list --state open --json number,title,labels,assignees,createdAt > open_issues.json\n" +
			"      gh pr list --state open --json number,title,author,createdAt > open_prs.json\n"
	case "dependency-monitor":
		return "      cat package.json go.mod requirements.txt 2>/dev/null > deps_snapshot.txt\n"
	case "upstream-monitor":
		return "      git log --oneline -20 > recent_commits.txt\n"
	default:
		return "      echo 'Gathering data...'\n"
	}
}

func promptBody(cfg WorkflowConfig) string {
	var b strings.Builder
	id := cfg.Archetype.ID

	switch id {
	case "issue-triage":
		b.WriteString("# Issue Triage\n\n")
		b.WriteString("You are an issue triage assistant. When a new issue is opened:\n\n")
		b.WriteString("1. Read the issue title and body carefully\n")
		b.WriteString("2. Apply appropriate labels based on content (bug, feature, question, documentation)\n")
		b.WriteString("3. If the issue lacks reproduction steps for a bug, ask the author to provide them\n")
		b.WriteString("4. Add a brief triage comment summarizing your assessment\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- Close issues without explanation\n")
		b.WriteString("- Apply more than 3 labels\n")
		b.WriteString("- Modify issue titles\n")
		b.WriteString("- Assign issues to people\n")

	case "status-report":
		b.WriteString("# Status Report Generator\n\n")
		b.WriteString("Generate a status report summarizing recent repository activity.\n\n")
		b.WriteString("Use the pre-fetched data files:\n")
		b.WriteString("- `open_issues.json` — current open issues\n")
		b.WriteString("- `open_prs.json` — current open pull requests\n\n")
		b.WriteString("## Report Format\n")
		b.WriteString("Create an issue with:\n")
		b.WriteString("- Summary of open issues by label/category\n")
		b.WriteString("- Open PR status and age\n")
		b.WriteString("- Key trends or items needing attention\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- Make any code changes\n")
		b.WriteString("- Close or modify existing issues\n")
		b.WriteString("- Create duplicate reports — check for existing ones first\n")

	case "code-improvement":
		b.WriteString("# Code Improvement\n\n")
		b.WriteString("Analyze the codebase and identify improvements:\n\n")
		b.WriteString("1. Scan for common issues:\n")
		b.WriteString("   - Performance problems (N+1 queries, unbounded loops, memory leaks)\n")
		b.WriteString("   - Security concerns (hardcoded secrets, SQL injection, XSS)\n")
		b.WriteString("   - Code quality (dead code, duplicated logic, missing error handling)\n")
		b.WriteString("2. For each finding, create a focused PR with the fix\n")
		b.WriteString("3. Explain what was wrong and why the fix is correct in the PR body\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- Reformat or restyle code not related to the fix\n")
		b.WriteString("- Combine unrelated fixes in one PR\n")
		b.WriteString("- Change public API signatures without explanation\n")

	case "pr-review":
		b.WriteString("# Pull Request Review\n\n")
		b.WriteString("Review the pull request for quality and correctness:\n\n")
		b.WriteString("1. Check for bugs, logic errors, and security issues\n")
		b.WriteString("2. Verify tests cover the changes\n")
		b.WriteString("3. Look for performance regressions\n")
		b.WriteString("4. Leave constructive inline comments on specific lines\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- Comment on style or formatting preferences\n")
		b.WriteString("- Approve or request changes — only leave comments\n")
		b.WriteString("- Rewrite the author's approach entirely\n")

	case "documentation-updater":
		b.WriteString("# Documentation Updater\n\n")
		b.WriteString("Keep documentation accurate and current:\n\n")
		b.WriteString("1. Compare docs to actual code behavior\n")
		b.WriteString("2. Update outdated examples, API references, and guides\n")
		b.WriteString("3. Add documentation for undocumented public APIs\n")
		b.WriteString("4. Create a PR with the updates\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- Change code — only update documentation files\n")
		b.WriteString("- Remove existing documentation without replacement\n")
		b.WriteString("- Add opinions or commentary — keep docs factual\n")

	case "upstream-monitor":
		b.WriteString("# Upstream Monitor\n\n")
		b.WriteString("Track upstream changes and sync when needed:\n\n")
		b.WriteString("1. Check for new releases or breaking changes in dependencies\n")
		b.WriteString("2. Create an issue summarizing what changed and potential impact\n")
		b.WriteString("3. If the change is straightforward, create a PR with the update\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- Auto-merge dependency updates\n")
		b.WriteString("- Update multiple dependencies in one PR\n")
		b.WriteString("- Ignore breaking changes — always flag them\n")

	case "dependency-monitor":
		b.WriteString("# Dependency Monitor\n\n")
		b.WriteString("Track and update project dependencies:\n\n")
		b.WriteString("1. Check for outdated dependencies\n")
		b.WriteString("2. Review changelogs for breaking changes\n")
		b.WriteString("3. Create individual PRs for updates with clear descriptions\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- Batch unrelated dependency updates\n")
		b.WriteString("- Update major versions without flagging breaking changes\n")
		b.WriteString("- Skip reading changelogs\n")

	case "content-moderation":
		b.WriteString("# Content Moderation\n\n")
		b.WriteString("Review content for quality and policy compliance:\n\n")
		b.WriteString("1. Check new issues and discussions for spam or policy violations\n")
		b.WriteString("2. Flag problematic content with a label\n")
		b.WriteString("3. Leave a comment explaining any moderation action\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- Delete or hide content — only flag it\n")
		b.WriteString("- Make subjective quality judgments\n")
		b.WriteString("- Moderate maintainer or contributor content\n")

	default: // custom
		b.WriteString("# Custom Workflow\n\n")
		b.WriteString("Describe your workflow's purpose and behavior here.\n\n")
		b.WriteString("## Steps\n")
		b.WriteString("1. ...\n")
		b.WriteString("2. ...\n\n")
		b.WriteString("## DO NOT\n")
		b.WriteString("- ...\n")
	}

	// Project context
	if cfg.ProjectContext != "" {
		b.WriteString("\n\n## Project Context\n\n")
		b.WriteString(cfg.ProjectContext)
		b.WriteString("\n")
	}

	return b.String()
}

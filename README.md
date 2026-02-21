# gh-aw-create

Interactive terminal wizard for creating [GitHub Agentic Workflows](https://github.github.com/gh-aw/). Same workflow generator as [ashleywolf.github.io/agentic-prompt-generator](https://ashleywolf.github.io/agentic-prompt-generator/), but in your terminal.

## Demo

```
$ gh aw-create

  ✓ Type  ─  ● Triggers  ─  ○ Context  ─  ○ Generate

  Triggers for Status Report
  Recommended triggers are pre-selected — adjust as needed

  ▸ [✓] schedule       Run on a cron schedule
    [✓] workflow_dispatch  Manual trigger from Actions tab
    [✓] issues         When issues are opened or edited
    [ ] pull_request   When PRs are opened or updated
    ...

  ↑↓ navigate • space toggle • enter next • esc back
```

## Install

```bash
gh extension install ashleywolf/gh-aw-create
```

### Build from source

```bash
git clone https://github.com/ashleywolf/gh-aw-create.git
cd gh-aw-create
go build -o bin/gh-aw-create .
```

## Usage

```bash
gh aw-create
```

The wizard walks you through:

1. **Pick a workflow type** — issue triage, status reports, dependency monitoring, code improvement, and more
2. **Select triggers** — recommended triggers are pre-selected based on your workflow type
3. **Add context** — optional project details and memory toggle
4. **Preview & save** — review the generated `.md` file and write it to `.github/workflows/`

After saving, the tool shows next steps:

```
✓ Written to .github/workflows/status-report.md

Next steps

  1. Ensure GitHub Actions is enabled on your repo
  2. gh extension install github/gh-aw
  3. gh aw add-wizard
  4. gh aw compile .github/workflows/status-report.md
  5. git add .github/workflows/status-report.md .github/workflows/status-report.lock.yml
  6. git commit -m 'Add agentic workflow' && git push
  7. gh aw run Status Report
```

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| `↑` / `↓` / `j` / `k` | Navigate |
| `enter` | Select / next step |
| `space` / `x` | Toggle trigger |
| `tab` | Toggle memory |
| `w` | Write file to `.github/workflows/` |
| `esc` | Go back |
| `q` / `ctrl+c` | Quit |

## How it works

The wizard embeds the same `patterns.json` data as the [web generator](https://github.com/ashleywolf/agentic-prompt-generator). It auto-infers capabilities (pre-steps, bash, GitHub toolsets) based on your workflow type — you don't need to know the internals of gh-aw frontmatter.

## Requirements

- [GitHub CLI](https://cli.github.com/) (`gh`)
- [gh-aw extension](https://github.com/github/gh-aw) (for compiling and running workflows)

## License

MIT

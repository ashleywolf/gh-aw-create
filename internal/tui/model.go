package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ashleywolf/gh-aw-create/internal/data"
	"github.com/ashleywolf/gh-aw-create/internal/generator"
)

type step int

const (
	stepArchetype step = iota
	stepTriggers
	stepContext
	stepPreview
)

var stepLabels = []string{"Type", "Triggers", "Context", "Generate"}

type Model struct {
	patterns  *data.Patterns
	step      step
	width     int
	height    int
	quitting  bool
	generated string
	written   bool
	writePath string
	writeErr  string

	// Step 1: archetype
	archCursor int

	// Step 2: triggers
	triggerSelected map[string]bool
	triggerCursor   int

	// Step 3: context
	contextInput textinput.Model
	useMemory    bool

	// Step 4: preview
	previewScroll int
}

func NewModel(p *data.Patterns) Model {
	ti := textinput.New()
	ti.Placeholder = "e.g., Monorepo with packages in /packages/*, uses conventional commits..."
	ti.CharLimit = 500
	ti.Width = 70

	return Model{
		patterns:        p,
		triggerSelected: make(map[string]bool),
		contextInput:    ti,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.step == stepPreview {
				m.quitting = true
				return m, tea.Quit
			}
			if msg.String() == "ctrl+c" {
				m.quitting = true
				return m, tea.Quit
			}
		case "esc":
			if m.step > stepArchetype {
				m.step--
				return m, nil
			}
		}

		switch m.step {
		case stepArchetype:
			return m.updateArchetype(msg)
		case stepTriggers:
			return m.updateTriggers(msg)
		case stepContext:
			return m.updateContext(msg)
		case stepPreview:
			return m.updatePreview(msg)
		}
	}

	if m.step == stepContext {
		var cmd tea.Cmd
		m.contextInput, cmd = m.contextInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// --- Step 1: Archetype ---

func (m Model) updateArchetype(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.archCursor > 0 {
			m.archCursor--
		}
	case "down", "j":
		if m.archCursor < len(m.patterns.Archetypes)-1 {
			m.archCursor++
		}
	case "enter":
		// Pre-select recommended triggers
		arch := m.patterns.Archetypes[m.archCursor]
		m.triggerSelected = make(map[string]bool)
		for _, t := range arch.RecommendedTriggers {
			m.triggerSelected[t.Type] = true
		}
		m.triggerCursor = 0
		m.step = stepTriggers
	}
	return m, nil
}

func (m Model) viewArchetype() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render("What type of workflow?"))
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render("Pick the workflow that fits your use case"))
	b.WriteString("\n\n")

	for i, a := range m.patterns.Archetypes {
		emoji := data.ArchetypeEmoji(a.ID)
		cursor := "  "
		style := UnselectedItem
		descStyle := ItemDesc

		if i == m.archCursor {
			cursor = "▸ "
			style = SelectedItem
			descStyle = descStyle.Foreground(lipgloss.Color("#58a6ff"))
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, emoji, style.Render(a.Label)))
		b.WriteString(descStyle.Render(a.Description))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑↓ navigate • enter select • ctrl+c quit"))
	return b.String()
}

// --- Step 2: Triggers ---

func (m Model) updateTriggers(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.triggerCursor > 0 {
			m.triggerCursor--
		}
	case "down", "j":
		if m.triggerCursor < len(data.AllTriggers)-1 {
			m.triggerCursor++
		}
	case " ", "x":
		t := data.AllTriggers[m.triggerCursor]
		m.triggerSelected[t] = !m.triggerSelected[t]
	case "enter":
		m.contextInput.Focus()
		m.step = stepContext
		return m, m.contextInput.Focus()
	}
	return m, nil
}

func (m Model) viewTriggers() string {
	var b strings.Builder
	arch := m.patterns.Archetypes[m.archCursor]
	b.WriteString(TitleStyle.Render(fmt.Sprintf("Triggers for %s", arch.Label)))
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render("Recommended triggers are pre-selected — adjust as needed"))
	b.WriteString("\n\n")

	for i, t := range data.AllTriggers {
		cursor := "  "
		if i == m.triggerCursor {
			cursor = "▸ "
		}

		check := Unchecked.Render("[ ]")
		nameStyle := UnselectedItem
		if m.triggerSelected[t] {
			check = Checked.Render("[✓]")
			nameStyle = SelectedItem
		}

		desc := data.TriggerDescriptions[t]
		b.WriteString(fmt.Sprintf("%s%s %s  %s\n",
			cursor, check, nameStyle.Render(t), ItemDesc.Render(desc)))
	}

	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("↑↓ navigate • space toggle • enter next • esc back"))
	return b.String()
}

// --- Step 3: Context ---

func (m Model) updateContext(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.useMemory = !m.useMemory
		return m, nil
	case "enter":
		m.generateWorkflow()
		m.step = stepPreview
		return m, nil
	}
	var cmd tea.Cmd
	m.contextInput, cmd = m.contextInput.Update(msg)
	return m, cmd
}

func (m Model) viewContext() string {
	var b strings.Builder
	arch := m.patterns.Archetypes[m.archCursor]
	b.WriteString(TitleStyle.Render(fmt.Sprintf("Context for %s", arch.Label)))
	b.WriteString("\n")
	b.WriteString(SubtitleStyle.Render("Optional — add project details for a better workflow"))
	b.WriteString("\n\n")

	memCheck := Unchecked.Render("[ ]")
	memStyle := UnselectedItem
	if m.useMemory {
		memCheck = Checked.Render("[✓]")
		memStyle = SelectedItem
	}
	b.WriteString(fmt.Sprintf("  %s %s\n", memCheck, memStyle.Render("🧠 Remember across runs")))
	b.WriteString(ItemDesc.Render("Track trends and context between executions"))
	b.WriteString("\n\n")

	b.WriteString("  Project context:\n")
	b.WriteString("  " + m.contextInput.View())
	b.WriteString("\n\n")

	b.WriteString(HelpStyle.Render("tab toggle memory • enter generate • esc back"))
	return b.String()
}

// --- Step 4: Preview ---

func (m *Model) generateWorkflow() {
	arch := m.patterns.Archetypes[m.archCursor]
	var triggers []string
	for _, t := range data.AllTriggers {
		if m.triggerSelected[t] {
			triggers = append(triggers, t)
		}
	}

	cfg := generator.WorkflowConfig{
		Archetype:      arch,
		Triggers:       triggers,
		ProjectContext: m.contextInput.Value(),
		UseMemory:      m.useMemory,
	}
	m.generated = generator.Generate(cfg)
}

func (m Model) updatePreview(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.previewScroll > 0 {
			m.previewScroll--
		}
	case "down", "j":
		m.previewScroll++
	case "w":
		m.writeFile()
	case "q":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) writeFile() {
	if m.written {
		return
	}
	arch := m.patterns.Archetypes[m.archCursor]
	name := strings.ReplaceAll(strings.ToLower(arch.Label), " ", "-")
	m.writePath = fmt.Sprintf(".github/workflows/%s.md", name)

	// Create directory and write file
	err := writeWorkflowFile(m.writePath, m.generated)
	if err != nil {
		m.writeErr = err.Error()
		return
	}
	m.written = true
}

func (m Model) viewPreview() string {
	var b strings.Builder
	arch := m.patterns.Archetypes[m.archCursor]

	if m.written {
		b.WriteString(SuccessStyle.Render(fmt.Sprintf("✓ Written to %s", m.writePath)))
		b.WriteString("\n\n")
		b.WriteString(TitleStyle.Render("Next steps"))
		b.WriteString("\n\n")
		steps := []struct{ num, cmd, desc string }{
			{"1", "", "Ensure GitHub Actions is enabled on your repo"},
			{"2", "gh extension install github/gh-aw", "Install the gh-aw extension (if not already)"},
			{"3", "gh aw add-wizard", "Set up your AI engine secret"},
			{"4", fmt.Sprintf("gh aw compile %s", m.writePath), "Compile the workflow"},
			{"5", fmt.Sprintf("git add %s %s", m.writePath, strings.Replace(m.writePath, ".md", ".lock.yml", 1)), "Stage both files"},
			{"6", "git commit -m 'Add agentic workflow' && git push", "Commit and push"},
			{"7", fmt.Sprintf("gh aw run %s", strings.TrimSuffix(arch.Label, " ")), "Trigger your first run"},
		}
		for _, s := range steps {
			if s.cmd != "" {
				b.WriteString(fmt.Sprintf("  %s. %s\n", s.num, NextStepCmd.Render(s.cmd)))
				b.WriteString(fmt.Sprintf("     %s\n", NextStepStyle.Render(s.desc)))
			} else {
				b.WriteString(fmt.Sprintf("  %s. %s\n", s.num, NextStepStyle.Render(s.desc)))
			}
		}
		b.WriteString("\n")
		b.WriteString(HelpStyle.Render("q quit"))
		return b.String()
	}

	if m.writeErr != "" {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#f85149")).Render("Error: " + m.writeErr))
		b.WriteString("\n\n")
	}

	b.WriteString(TitleStyle.Render(fmt.Sprintf("Preview: %s workflow", arch.Label)))
	b.WriteString("\n")

	// Show scrollable preview
	lines := strings.Split(m.generated, "\n")
	maxLines := m.height - 10
	if maxLines < 5 {
		maxLines = 20
	}
	if m.previewScroll >= len(lines) {
		m.previewScroll = len(lines) - 1
	}
	end := m.previewScroll + maxLines
	if end > len(lines) {
		end = len(lines)
	}
	visible := strings.Join(lines[m.previewScroll:end], "\n")
	b.WriteString(PreviewBox.Render(visible))
	b.WriteString("\n\n")

	b.WriteString(HelpStyle.Render("↑↓ scroll • w write to .github/workflows/ • esc back • q quit"))
	return b.String()
}

// --- View ---

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var content string
	switch m.step {
	case stepArchetype:
		content = m.viewArchetype()
	case stepTriggers:
		content = m.viewTriggers()
	case stepContext:
		content = m.viewContext()
	case stepPreview:
		content = m.viewPreview()
	}

	progress := m.renderProgress()
	return AppStyle.Render(progress + "\n\n" + content)
}

func (m Model) renderProgress() string {
	var parts []string
	for i, label := range stepLabels {
		s := step(i)
		var rendered string
		switch {
		case s < m.step:
			rendered = ProgressDone.Render(fmt.Sprintf("✓ %s", label))
		case s == m.step:
			rendered = ProgressActive.Render(fmt.Sprintf("● %s", label))
		default:
			rendered = ProgressPending.Render(fmt.Sprintf("○ %s", label))
		}
		parts = append(parts, rendered)
	}
	return strings.Join(parts, ProgressPending.Render("  ─  "))
}

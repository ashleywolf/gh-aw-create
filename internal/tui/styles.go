package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	blue      = lipgloss.Color("#58a6ff")
	green     = lipgloss.Color("#3fb950")
	purple    = lipgloss.Color("#bc8cff")
	dimWhite  = lipgloss.Color("#8b949e")
	white     = lipgloss.Color("#e6edf3")
	darkBg    = lipgloss.Color("#161b22")
	cardBg    = lipgloss.Color("#0d1117")
	borderDim = lipgloss.Color("#30363d")

	// App frame
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	// Title
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(blue).
			MarginBottom(1)

	// Subtitle / step description
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(dimWhite).
			MarginBottom(1)

	// Progress bar
	ProgressActive = lipgloss.NewStyle().
			Foreground(blue).
			Bold(true)

	ProgressDone = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	ProgressPending = lipgloss.NewStyle().
			Foreground(dimWhite)

	// List items
	SelectedItem = lipgloss.NewStyle().
			Foreground(blue).
			Bold(true)

	UnselectedItem = lipgloss.NewStyle().
			Foreground(white)

	ItemDesc = lipgloss.NewStyle().
			Foreground(dimWhite).
			PaddingLeft(4)

	// Checked/unchecked
	Checked   = lipgloss.NewStyle().Foreground(green).Bold(true)
	Unchecked = lipgloss.NewStyle().Foreground(dimWhite)

	// Preview box
	PreviewBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderDim).
			Padding(1, 2).
			MarginTop(1)

	// Help / footer
	HelpStyle = lipgloss.NewStyle().
			Foreground(dimWhite).
			MarginTop(1)

	// Success message
	SuccessStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	// Next steps
	NextStepStyle = lipgloss.NewStyle().
			Foreground(white).
			PaddingLeft(2)

	NextStepCmd = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true)
)

package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ashleywolf/gh-aw-create/internal/data"
	"github.com/ashleywolf/gh-aw-create/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "gh-aw-create",
	Short: "Create GitHub Agentic Workflows from the terminal",
	Long:  "Interactive TUI wizard for generating production-ready .md workflow files for GitHub Agentic Workflows (gh-aw).",
	RunE: func(cmd *cobra.Command, args []string) error {
		patterns, err := data.LoadPatterns()
		if err != nil {
			return fmt.Errorf("loading patterns: %w", err)
		}

		m := tui.NewModel(patterns)
		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

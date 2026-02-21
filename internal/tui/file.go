package tui

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeWorkflowFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	return os.WriteFile(path, []byte(content), 0644)
}

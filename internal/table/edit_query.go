package table

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) editAndRerunQuery() (tea.Model, tea.Cmd) {
	// Open editor with current query
	editedQuery, err := m.openQueryEditor()
	if err != nil {
		// Don't re-enter alt screen on error, just stay out
		fmt. Fprintf(os.Stderr, "\nError editing query: %v\n", err)
		fmt.Fprintf(os.Stderr, "Press any key to continue.. .\n")
		var buf [1]byte
		os. Stdin.Read(buf[:])
		// Return and quit - don't go back to table
		return m, tea.Quit
	}

	// Store the edited query and signal that we should re-run
	m.editedQuery = editedQuery
	m.shouldRerunQuery = true

	// Quit the TUI - the calling code will handle re-execution
	return m, tea.Quit
}

func (m Model) openQueryEditor() (string, error) {
	editorCmd := os. Getenv("EDITOR")
	if editorCmd == "" {
		editorCmd = "vim"
	}

	tmpFile, err := os.CreateTemp("", "pam-edit-query-*.sql")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	// Write current query to temp file
	if _, err := tmpFile.WriteString(m.currentQuery.SQL); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Exit alt screen before opening editor
	tea.ExitAltScreen()

	// Open editor
	cmd := exec.Command(editorCmd, tmpPath)
	cmd.Stdin = os. Stdin
	cmd.Stdout = os.Stdout
	cmd. Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Don't re-enter on error
		return "", fmt.Errorf("failed to run editor: %v", err)
	}

	// Give the terminal a moment to settle after editor closes
	time.Sleep(50 * time.Millisecond)

	// Reset terminal using stty (Unix/Linux/macOS)
	resetCmd := exec.Command("stty", "sane")
	resetCmd.Stdin = os.Stdin
	resetCmd.Stdout = os. Stdout
	resetCmd. Stderr = os.Stderr
	resetCmd.Run() // Ignore errors - this is best-effort

	// Read edited content
	editedData, err := os.ReadFile(tmpPath)
	if err != nil {
		// Don't re-enter on error
		return "", fmt.Errorf("failed to read edited file: %v", err)
	}

	editedSQL := strings.TrimSpace(string(editedData))
	if editedSQL == "" {
		// Don't re-enter on error
		return "", fmt.Errorf("empty query")
	}

	return editedSQL, nil
}

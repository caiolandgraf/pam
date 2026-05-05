package table

import (
	"os"
	"time"

	"github.com/caiolandgraf/pam/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) editAndRerunQuery() (tea.Model, tea.Cmd) {
	editorCmd := os.Getenv("EDITOR")
	if editorCmd == "" {
		editorCmd = "vim"
	}

	tmpFile, err := os.CreateTemp("", "pam-edit-query-*.sql")
	if err != nil {
		return m, nil
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.WriteString(m.currentQuery.SQL); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return m, nil
	}
	tmpFile.Close()

	// Verify temp file exists before opening editor
	if _, err := os.Stat(tmpPath); err != nil {
		os.Remove(tmpPath)
		return m, nil
	}

	return m.openInlineEditor(
		editorKindEditQuery,
		"Edit query",
		m.currentQuery.SQL,
		m.selectedCol,
	)
}

type queryEditCompleteMsg struct {
	sql       string
	cancelled bool
}

func (m Model) handleQueryEditComplete(
	msg queryEditCompleteMsg,
) (tea.Model, tea.Cmd) {
	// If user cancelled (exited without saving)
	if msg.cancelled {
		m.statusMessage = styles.Error.Render("✗ Edit canceled")
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return blinkMsg{}
		})
	}

	m.editedQuery = msg.sql
	m.shouldRerunQuery = true

	return m, tea.Quit
}

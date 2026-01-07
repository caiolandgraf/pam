package table

import (
	"os"
	"strings"

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

	if _, err := tmpFile.WriteString(m.currentQuery. SQL); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return m, nil
	}
	tmpFile.Close()

	cmd := buildEditorCommand(editorCmd, tmpPath, m.currentQuery. SQL, CursorAtEndOfFile)
	
	return m, tea. ExecProcess(cmd, func(err error) tea.Msg {
		editedData, readErr := os.ReadFile(tmpPath)
		os. Remove(tmpPath)
		
		if err != nil || readErr != nil {
			return nil
		}

		editedSQL := strings.TrimSpace(string(editedData))
		if editedSQL == "" {
			return nil
		}

		return queryEditCompleteMsg{
			sql:  editedSQL,
		}
	})
}

type queryEditCompleteMsg struct {
	sql string
}

func (m Model) handleQueryEditComplete(msg queryEditCompleteMsg) (tea.Model, tea.Cmd) {
	m.editedQuery = msg.sql
	m.shouldRerunQuery = true

	return m, tea.Quit
}

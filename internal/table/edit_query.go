package table

import (
	"strings"
	"time"

	"github.com/caiolandgraf/pam/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) editAndRerunQuery() (tea.Model, tea.Cmd) {
	if strings.TrimSpace(m.currentQuery.SQL) == "" {
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

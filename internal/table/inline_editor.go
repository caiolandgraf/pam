package table

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/caiolandgraf/pam/internal/styles"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) openInlineEditor(
	kind editorKind,
	title string,
	content string,
	colIndex int,
) (Model, tea.Cmd) {
	ta := textarea.New()
	ta.SetValue(content)
	ta.Focus()
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.Prompt = "  "

	width := m.width - 6
	if width < 40 {
		width = 40
	}
	height := m.height - 10
	if height < 6 {
		height = 6
	}
	ta.SetWidth(width)
	ta.SetHeight(height)

	m.editorActive = true
	m.editorKind = kind
	m.editorTitle = title
	m.editorHelp = "Ctrl+S: save • Esc: cancel"
	m.editor = ta
	m.editorCol = colIndex

	return m, nil
}

func (m Model) handleInlineEditorUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.editorActive = false
			m.editorHelp = ""
			m.statusMessage = styles.Error.Render("✗ Edit canceled")
			return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
				return blinkMsg{}
			})
		case "ctrl+s":
			return m.saveInlineEditor()
		}
	}

	var cmd tea.Cmd
	m.editor, cmd = m.editor.Update(msg)
	return m, cmd
}

func (m Model) saveInlineEditor() (tea.Model, tea.Cmd) {
	content := strings.TrimSpace(m.editor.Value())

	m.editorActive = false
	m.editorHelp = ""

	if content == "" {
		m.statusMessage = styles.Error.Render("✗ Edit canceled")
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return blinkMsg{}
		})
	}

	switch m.editorKind {
	case editorKindUpdateCell:
		return m.handleEditorComplete(editorCompleteMsg{
			sql:       content,
			colIndex:  m.editorCol,
			cancelled: false,
		})
	case editorKindEditQuery:
		return m.handleQueryEditComplete(queryEditCompleteMsg{
			sql:       content,
			cancelled: false,
		})
	case editorKindDetailUpdate:
		return m.handleDetailViewEditComplete(detailViewEditCompleteMsg{
			sql:      content,
			colIndex: m.editorCol,
		})
	default:
		m.statusMessage = styles.Error.Render("✗ Unknown editor mode")
		return m, nil
	}
}

func (m Model) renderInlineEditorView() string {
	var b strings.Builder

	title := "✎ " + m.editorTitle
	b.WriteString(styles.Title.Render(title))
	b.WriteString("\n")

	sepWidth := m.width
	if sepWidth < 30 {
		sepWidth = 30
	}
	b.WriteString(styles.Separator.Render(strings.Repeat("─", sepWidth)))
	b.WriteString("\n")

	b.WriteString(m.editor.View())
	b.WriteString("\n")

	b.WriteString(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
			Render(m.editorHelp),
	)

	return b.String()
}

func (m Model) requestDeleteRows(rows []int) (tea.Model, tea.Cmd) {
	if len(rows) == 0 {
		return m, nil
	}

	m.confirmActive = true
	m.confirmRows = normalizeRows(rows)
	m.confirmMessage = fmt.Sprintf(
		"Delete %d row(s)? (y/N)",
		len(m.confirmRows),
	)
	return m, nil
}

func (m Model) handleDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		return m.executeDeleteRows(m.confirmRows)
	case "n", "N", "esc":
		m.confirmActive = false
		m.confirmMessage = ""
		m.confirmRows = nil
		return m, nil
	default:
		return m, nil
	}
}

func (m Model) executeDeleteRows(rows []int) (tea.Model, tea.Cmd) {
	m.confirmActive = false
	m.confirmMessage = ""

	if m.primaryKeyCol == "" {
		m.statusMessage = styles.Error.Render("✗ Delete requires a primary key")
		return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
			return blinkMsg{}
		})
	}

	rows = normalizeRows(rows)
	deleted := 0

	for _, row := range rows {
		if row < 0 || row >= m.numRows() {
			continue
		}
		m.selectedRow = row

		stmt := m.buildDeleteStatement()
		if err := validateDeleteStatement(stmt); err != nil {
			m.statusMessage = styles.Error.Render(
				fmt.Sprintf("✗ Delete failed: %v", err),
			)
			return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
				return blinkMsg{}
			})
		}
		if err := m.executeDelete(stmt); err != nil {
			m.statusMessage = styles.Error.Render(
				fmt.Sprintf("✗ Delete failed: %v", err),
			)
			return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
				return blinkMsg{}
			})
		}

		m.data = append(m.data[:row], m.data[row+1:]...)
		deleted++
	}

	m.clearMarkedRows()
	m.visualMode = false
	if m.selectedRow >= m.numRows() && m.numRows() > 0 {
		m.selectedRow = m.numRows() - 1
	}
	if m.offsetY >= m.numRows() && m.numRows() > 0 {
		m.offsetY = m.numRows() - 1
	}

	if deleted > 0 {
		m.statusMessage = styles.Success.Render(
			fmt.Sprintf("✓ Deleted %d row(s)", deleted),
		)
	}

	return m, tea.Batch(tea.ClearScreen, m.blinkCmd())
}

func normalizeRows(rows []int) []int {
	seen := map[int]struct{}{}
	for _, r := range rows {
		seen[r] = struct{}{}
	}

	out := make([]int, 0, len(seen))
	for r := range seen {
		out = append(out, r)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i] > out[j] // delete from bottom to top
	})

	return out
}

func (m Model) openValueEditor(
	kind editorKind,
	title string,
	currentValue string,
	colIndex int,
) (Model, tea.Cmd) {
	columnName := ""
	if colIndex >= 0 && colIndex < len(m.columns) {
		columnName = m.columns[colIndex]
	}
	masked := isSensitiveColumn(columnName)

	input := textinput.New()
	input.Prompt = "  New: "
	input.CharLimit = 0
	input.SetValue("")
	input.Placeholder = "type new value…"
	input.PlaceholderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Muted))
	if masked {
		input.EchoMode = textinput.EchoPassword
		input.EchoCharacter = '•'
	}
	input.Focus()

	width := m.width - 20
	if width < 20 {
		width = 20
	}
	input.Width = width

	m.valueEditorActive = true
	m.valueEditorKind = kind
	m.valueEditorTitle = title
	m.valueEditorHelp = "Enter/Tab: save • Esc: cancel"
	m.valueEditorInput = input
	m.valueEditorCurrent = currentValue
	m.valueEditorCol = colIndex
	m.valueEditorMasked = masked
	m.valueEditorPlaceholder = input.Placeholder

	return m, nil
}

func (m Model) handleValueEditorUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.valueEditorActive = false
			m.valueEditorHelp = ""
			m.statusMessage = styles.Error.Render("✗ Edit canceled")
			return m, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
				return blinkMsg{}
			})
		case "enter", "tab":
			return m.saveValueEditor()
		}
	}

	var cmd tea.Cmd
	m.valueEditorInput, cmd = m.valueEditorInput.Update(msg)
	return m, cmd
}

func (m Model) saveValueEditor() (tea.Model, tea.Cmd) {
	newValue := m.valueEditorInput.Value()

	m.valueEditorActive = false
	m.valueEditorHelp = ""

	updateSQL := m.buildUpdateStatementWithValue(newValue)

	switch m.valueEditorKind {
	case editorKindUpdateCell:
		return m.handleEditorComplete(editorCompleteMsg{
			sql:       updateSQL,
			colIndex:  m.valueEditorCol,
			cancelled: false,
		})
	case editorKindDetailUpdate:
		return m.handleDetailViewEditComplete(detailViewEditCompleteMsg{
			sql:      updateSQL,
			colIndex: m.valueEditorCol,
		})
	default:
		m.statusMessage = styles.Error.Render("✗ Unknown editor mode")
		return m, nil
	}
}

func (m Model) renderValueEditorView() string {
	var b strings.Builder

	title := "✎ " + m.valueEditorTitle
	b.WriteString(styles.Title.Render(title))
	b.WriteString("\n")

	sepWidth := m.width
	if sepWidth < 30 {
		sepWidth = 30
	}
	b.WriteString(styles.Separator.Render(strings.Repeat("─", sepWidth)))
	b.WriteString("\n")

	displayCurrent := m.valueEditorCurrent
	if m.valueEditorMasked && displayCurrent != "" {
		displayCurrent = strings.Repeat("•", len(displayCurrent))
	}

	b.WriteString(styles.Faint.Render("Current: "))
	b.WriteString(styles.TableCell.Render(displayCurrent))
	b.WriteString("\n\n")

	b.WriteString(m.valueEditorInput.View())
	b.WriteString("\n\n")

	b.WriteString(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
			Render(m.valueEditorHelp),
	)

	return b.String()
}

func isSensitiveColumn(name string) bool {
	n := strings.ToLower(strings.TrimSpace(name))
	if n == "" {
		return false
	}

	if strings.Contains(n, "password") ||
		strings.Contains(n, "passwd") ||
		strings.Contains(n, "pwd") ||
		strings.Contains(n, "secret") ||
		strings.Contains(n, "token") ||
		strings.Contains(n, "apikey") ||
		strings.Contains(n, "api_key") {
		return true
	}

	return strings.HasSuffix(n, "_key")
}

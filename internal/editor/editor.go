package editor

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/eduardofuncao/pam/internal/db"
	"github.com/eduardofuncao/pam/internal/styles"
)

type EditorModel struct {
	textArea  textarea.Model
	width     int
	height    int
	query     db.Query
	submitted bool
}

func NewEditor(initialQuery db.Query) EditorModel {
	ta := textarea.New()
	ta.Placeholder = "Enter your query..."
	ta.Focus()
	ta.CharLimit = 10000
	ta.SetWidth(80)

	formattedSQL := FormatSQLWithLineBreaks(initialQuery.SQL)
	ta.SetValue(formattedSQL)

	lineCount := countLines(formattedSQL)
	height := min(max(lineCount, 3), 15) // Min 3 lines, max 15 lines
	ta.SetHeight(height)

	return EditorModel{
		textArea:  ta,
		query:     initialQuery,
		submitted: false,
	}
}

func (m EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlD: // Use Ctrl+D for execution
			m.query.SQL = strings.TrimSpace(m.textArea.Value())
			m.submitted = true
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textArea.SetWidth(min(m.width-4, 120))
	}

	m.textArea, cmd = m.textArea.Update(msg)

	lineCount := countLines(m.textArea.Value())
	newHeight := min(max(lineCount, 3), 15)
	if newHeight != m.textArea.Height() {
		m.textArea.SetHeight(newHeight)
	}

	return m, cmd
}

func (m EditorModel) View() string {
	var content string
	if m.submitted {
		highlightedSQL := HighlightSQL(m.textArea.Value())
		content = fmt.Sprintf(
			"%s\n%s\n%s",
			styles.Title.Render("\n◆ "+m.query.Name),
			highlightedSQL,
			styles.Separator.Render("──────────────────────────────────────────────────────────"),
		)
	} else {
		content = fmt.Sprintf(
			"%s\n%s\n%s\n%s",
			styles.Title.Render("\n◆ "+m.query.Name),
			m.textArea.View(),
			styles.Faint.Render("Ctrl+D: Execute Query | Esc/Ctrl+C: Cancel"),
			styles.Separator.Render("──────────────────────────────────────────────────────────"),
		)
	}

	return content + "\n"
}

func (m EditorModel) GetQuery() (db.Query, bool) {
	return m.query, m.submitted
}

func EditQuery(query db.Query, edit bool) (db.Query, bool, error) {
	if !edit {
		m := NewEditor(query)
		m.submitted = true
		fmt.Print(m.View())
		return query, false, nil
	}

	m := NewEditor(query)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return db.Query{}, false, err
	}

	editor := finalModel.(EditorModel)
	query, submitted := editor.GetQuery()
	return query, submitted, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

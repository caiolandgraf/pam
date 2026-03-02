package table

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caiolandgraf/pam/internal/db"
	"github.com/caiolandgraf/pam/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableViewModel is the bubbletea model for the table structure view
type TableViewModel struct {
	width        int
	height       int
	tableName    string
	columns      []db.ColumnInfo
	selectedRow  int
	offsetY      int
	visibleRows  int
	conn         db.DatabaseConnection
	elapsed      time.Duration
	message      string // transient status message
	messageStyle lipgloss.Style
}

// TableViewResult holds the outcome of a table-view session
type TableViewResult struct {
	ShouldRefresh bool
}

// NewTableViewModel creates a new table structure viewer
func NewTableViewModel(
	tableName string,
	columns []db.ColumnInfo,
	conn db.DatabaseConnection,
	elapsed time.Duration,
) TableViewModel {
	return TableViewModel{
		tableName:    tableName,
		columns:      columns,
		selectedRow:  0,
		offsetY:      0,
		conn:         conn,
		elapsed:      elapsed,
		messageStyle: styles.Success,
	}
}

func (m TableViewModel) Init() tea.Cmd {
	return nil
}

// ---- Update ----

type tableViewBlinkMsg struct{}

type tableViewEditorMsg struct {
	sql string
	err error
}

func (m TableViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		reserved := 12 // header + footer + separators
		m.visibleRows = m.height - reserved
		if m.visibleRows < 3 {
			m.visibleRows = 3
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case tableViewBlinkMsg:
		m.message = ""
		return m, nil

	case tableViewEditorMsg:
		return m.handleEditorResult(msg)
	}

	return m, nil
}

func (m TableViewModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.selectedRow > 0 {
			m.selectedRow--
			if m.selectedRow < m.offsetY {
				m.offsetY = m.selectedRow
			}
		}
		return m, nil

	case "down", "j":
		if m.selectedRow < len(m.columns)-1 {
			m.selectedRow++
			if m.selectedRow >= m.offsetY+m.visibleRows {
				m.offsetY = m.selectedRow - m.visibleRows + 1
			}
		}
		return m, nil

	case "g":
		m.selectedRow = 0
		m.offsetY = 0
		return m, nil

	case "G":
		m.selectedRow = len(m.columns) - 1
		m.offsetY = m.selectedRow - m.visibleRows + 1
		if m.offsetY < 0 {
			m.offsetY = 0
		}
		return m, nil

	case "pgup", "ctrl+u":
		m.selectedRow -= m.visibleRows
		if m.selectedRow < 0 {
			m.selectedRow = 0
		}
		m.offsetY = m.selectedRow
		return m, nil

	case "pgdown", "ctrl+d":
		m.selectedRow += m.visibleRows
		if m.selectedRow >= len(m.columns) {
			m.selectedRow = len(m.columns) - 1
		}
		if m.selectedRow >= m.offsetY+m.visibleRows {
			m.offsetY = m.selectedRow - m.visibleRows + 1
		}
		return m, nil

	case "a":
		// Add column
		return m.addColumn()

	case "e":
		// Edit selected column (alter)
		return m.editColumn()

	case "r":
		// Rename selected column
		return m.renameColumn()

	case "D":
		// Drop selected column
		return m.dropColumn()

	case "enter":
		// Show detail of selected column
		return m, nil
	}

	return m, nil
}

func (m TableViewModel) handleEditorResult(
	msg tableViewEditorMsg,
) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.message = fmt.Sprintf("Editor error: %v", msg.err)
		m.messageStyle = styles.Error
		return m, m.blinkCmd()
	}

	sql := strings.TrimSpace(msg.sql)
	if sql == "" {
		m.message = "Cancelled (empty SQL)"
		m.messageStyle = styles.Faint
		return m, m.blinkCmd()
	}

	// Execute the SQL
	if err := m.conn.Exec(sql); err != nil {
		m.message = fmt.Sprintf("✗ Error: %v", err)
		m.messageStyle = styles.Error
		return m, m.blinkCmd()
	}

	// Refresh column data
	newCols, err := m.conn.GetColumnDetails(m.tableName)
	if err != nil {
		m.message = fmt.Sprintf("✓ Executed, but refresh failed: %v", err)
		m.messageStyle = styles.Success
		return m, m.blinkCmd()
	}

	m.columns = newCols
	if m.selectedRow >= len(m.columns) && len(m.columns) > 0 {
		m.selectedRow = len(m.columns) - 1
	}

	m.message = "✓ Column updated successfully"
	m.messageStyle = styles.Success
	return m, tea.Batch(tea.ClearScreen, m.blinkCmd())
}

func (m TableViewModel) blinkCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return tableViewBlinkMsg{}
	})
}

// ---- Actions ----

func (m TableViewModel) addColumn() (tea.Model, tea.Cmd) {
	template := m.conn.BuildAddColumnSQL(
		m.tableName,
		"new_column",
		"VARCHAR(255)",
		true,
		"",
	)

	header := fmt.Sprintf(
		"-- ADD COLUMN to %s\n-- Edit the statement below and save to execute.\n-- To cancel, delete all content and save.\n\n",
		m.tableName,
	)

	return m.openEditor(header + template)
}

func (m TableViewModel) editColumn() (tea.Model, tea.Cmd) {
	if len(m.columns) == 0 {
		return m, nil
	}

	col := m.columns[m.selectedRow]
	nullable := strings.ToUpper(col.Nullable) == "YES"
	defVal := col.DefaultValue
	if defVal == "NULL" {
		defVal = ""
	}

	template := m.conn.BuildAlterColumnSQL(
		m.tableName,
		col.Name,
		col.DataType,
		nullable,
		defVal,
	)

	header := fmt.Sprintf(
		"-- ALTER COLUMN '%s' on table '%s'\n-- Edit the statement below and save to execute.\n-- To cancel, delete all content and save.\n\n",
		col.Name,
		m.tableName,
	)

	return m.openEditor(header + template)
}

func (m TableViewModel) renameColumn() (tea.Model, tea.Cmd) {
	if len(m.columns) == 0 {
		return m, nil
	}

	col := m.columns[m.selectedRow]
	template := m.conn.BuildRenameColumnSQL(
		m.tableName,
		col.Name,
		col.Name+"_new",
	)

	header := fmt.Sprintf(
		"-- RENAME COLUMN '%s' on table '%s'\n-- Edit the new column name and save to execute.\n-- To cancel, delete all content and save.\n\n",
		col.Name,
		m.tableName,
	)

	return m.openEditor(header + template)
}

func (m TableViewModel) dropColumn() (tea.Model, tea.Cmd) {
	if len(m.columns) == 0 {
		return m, nil
	}

	col := m.columns[m.selectedRow]
	template := m.conn.BuildDropColumnSQL(m.tableName, col.Name)

	header := fmt.Sprintf(
		"-- DROP COLUMN '%s' from table '%s'\n-- WARNING: This will permanently remove the column and all its data!\n-- To cancel, delete all content and save.\n\n",
		col.Name,
		m.tableName,
	)

	return m.openEditor(header + template)
}

func (m TableViewModel) openEditor(content string) (tea.Model, tea.Cmd) {
	editorCmd := os.Getenv("EDITOR")
	if editorCmd == "" {
		editorCmd = "vim"
	}

	tmpFile, err := os.CreateTemp("", "pam-tableview-*.sql")
	if err != nil {
		m.message = fmt.Sprintf("Could not create temp file: %v", err)
		m.messageStyle = styles.Error
		return m, m.blinkCmd()
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		m.message = fmt.Sprintf("Could not write temp file: %v", err)
		m.messageStyle = styles.Error
		return m, m.blinkCmd()
	}
	tmpFile.Close()

	cmd := buildEditorCommand(editorCmd, tmpPath, content, CursorAtEndOfFile)

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		editedData, readErr := os.ReadFile(tmpPath)
		os.Remove(tmpPath)

		if err != nil {
			return tableViewEditorMsg{err: err}
		}
		if readErr != nil {
			return tableViewEditorMsg{err: readErr}
		}

		// Strip comment lines
		raw := string(editedData)
		var sqlLines []string
		for _, line := range strings.Split(raw, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
				sqlLines = append(sqlLines, trimmed)
			}
		}

		return tableViewEditorMsg{sql: strings.Join(sqlLines, " ")}
	})
}

// ---- View ----

func (m TableViewModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Title
	b.WriteString(styles.Title.Render(fmt.Sprintf("◆ Table: %s", m.tableName)))
	b.WriteString("\n")

	connInfo := fmt.Sprintf("%s/%s", m.conn.GetDbType(), m.conn.GetName())
	b.WriteString(styles.Faint.Render(connInfo))
	b.WriteString("\n\n")

	// Column headers
	colWidths := m.computeColumnWidths()
	headerCells := m.renderHeaderRow(colWidths)
	b.WriteString(headerCells)
	b.WriteString("\n")

	// Separator
	totalWidth := 0
	for i, w := range colWidths {
		totalWidth += w
		if i < len(colWidths)-1 {
			totalWidth += 1 // border
		}
	}
	b.WriteString(styles.Separator.Render(strings.Repeat("─", totalWidth)))
	b.WriteString("\n")

	// Data rows
	if len(m.columns) == 0 {
		b.WriteString(styles.Faint.Render("  No columns found"))
		b.WriteString("\n")
	} else {
		endRow := m.offsetY + m.visibleRows
		if endRow > len(m.columns) {
			endRow = len(m.columns)
		}

		for i := m.offsetY; i < endRow; i++ {
			b.WriteString(m.renderDataRow(i, colWidths))
			b.WriteString("\n")
		}
	}

	// Message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(m.messageStyle.Render(m.message))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString(m.renderFooter())

	return b.String()
}

func (m TableViewModel) computeColumnWidths() []int {
	// Columns: #  | Name | Type | Nullable | Default | PK | Extra
	headers := []string{
		"#",
		"Name",
		"Type",
		"Nullable",
		"Default",
		"PK",
		"Extra",
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, col := range m.columns {
		vals := m.columnToRow(col)
		for i, v := range vals {
			if len(v) > widths[i] {
				widths[i] = len(v)
			}
		}
	}

	// Apply min/max
	const minW = 3
	const maxW = 40
	for i := range widths {
		if widths[i] < minW {
			widths[i] = minW
		}
		if widths[i] > maxW {
			widths[i] = maxW
		}
		widths[i] += 2 // padding
	}

	// Check available width and scale down if needed
	totalBorders := (len(widths) - 1)
	totalNeeded := totalBorders
	for _, w := range widths {
		totalNeeded += w
	}

	available := m.width - 4
	if available < 40 {
		available = 40
	}

	if totalNeeded > available {
		// Shrink proportionally, but don't go below min
		scale := float64(
			available-totalBorders,
		) / float64(
			totalNeeded-totalBorders,
		)
		for i := range widths {
			scaled := int(float64(widths[i]) * scale)
			if scaled < minW+2 {
				scaled = minW + 2
			}
			widths[i] = scaled
		}
	}

	return widths
}

func (m TableViewModel) columnToRow(col db.ColumnInfo) []string {
	pkStr := ""
	if col.IsPrimaryKey {
		pkStr = "⚿ PK"
	}

	return []string{
		fmt.Sprintf("%d", col.OrdinalPos),
		col.Name,
		col.DataType,
		col.Nullable,
		col.DefaultValue,
		pkStr,
		col.Extra,
	}
}

func (m TableViewModel) renderHeaderRow(widths []int) string {
	headers := []string{
		"#",
		"Name",
		"Type",
		"Nullable",
		"Default",
		"PK",
		"Extra",
	}

	var cells []string
	for i, h := range headers {
		w := widths[i]
		content := tvFormatCell(h, w)
		cells = append(cells, styles.TableHeader.Render(content))
	}

	return strings.Join(cells, styles.TableBorder.Render("│"))
}

func (m TableViewModel) renderDataRow(rowIndex int, widths []int) string {
	col := m.columns[rowIndex]
	vals := m.columnToRow(col)

	var cells []string
	for i, v := range vals {
		w := widths[i]
		content := tvFormatCell(v, w)

		style := styles.TableCell
		if rowIndex == m.selectedRow {
			style = styles.TableSelected
		}

		// Highlight PK column name
		if i == 5 && col.IsPrimaryKey && rowIndex != m.selectedRow {
			style = styles.TableHeader
		}

		cells = append(cells, style.Render(content))
	}

	return strings.Join(cells, styles.TableBorder.Render("│"))
}

func (m TableViewModel) renderFooter() string {
	var b strings.Builder

	// Column count and timing
	info := fmt.Sprintf(
		"\n%s %s | %s",
		styles.Faint.Render(fmt.Sprintf("%d columns", len(m.columns))),
		styles.Faint.Render(fmt.Sprintf("In %.2fs", m.elapsed.Seconds())),
		styles.Faint.Render(
			fmt.Sprintf("[%d/%d]", m.selectedRow+1, len(m.columns)),
		),
	)
	b.WriteString(info)

	// Key hints
	add := styles.TableHeader.Render("a") + styles.Faint.Render("dd")
	edit := styles.TableHeader.Render("e") + styles.Faint.Render("dit")
	rename := styles.TableHeader.Render("r") + styles.Faint.Render("ename")
	drop := styles.TableHeader.Render("D") + styles.Faint.Render("rop")
	quit := styles.TableHeader.Render("q") + styles.Faint.Render("uit")
	hjkl := styles.TableHeader.Render("jk") + styles.Faint.Render("↓↑")

	b.WriteString(
		fmt.Sprintf(
			"  %s  %s  %s  %s  %s  %s",
			add,
			edit,
			rename,
			drop,
			quit,
			hjkl,
		),
	)

	return b.String()
}

func tvFormatCell(content string, width int) string {
	runes := []rune(content)
	runeCount := len(runes)

	if runeCount > width {
		return string(runes[:width-1]) + "…"
	}

	padding := width - runeCount
	return content + strings.Repeat(" ", padding)
}

// RenderTableView runs the table structure viewer TUI
func RenderTableView(
	tableName string,
	columns []db.ColumnInfo,
	conn db.DatabaseConnection,
	elapsed time.Duration,
) (TableViewModel, error) {
	model := NewTableViewModel(tableName, columns, conn, elapsed)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return model, err
	}
	return finalModel.(TableViewModel), nil
}

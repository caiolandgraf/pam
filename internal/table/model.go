package table

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eduardofuncao/pam/internal/db"
	"github.com/eduardofuncao/pam/internal/parser"
)

type Model struct {
	width             int
	height            int
	selectedRow       int
	selectedCol       int
	offsetX           int
	offsetY           int
	visibleCols       int
	visibleRows       int
	columns           []string
	columnTypes       []string
	data              [][]string
	elapsed           time.Duration
	blinkCopiedCell   bool
	visualMode        bool
	visualStartRow    int
	visualStartCol    int
	dbConnection      db.DatabaseConnection
	tableName         string
	primaryKeyCol     string
	blinkUpdatedCell  bool
	updatedRow        int
	updatedCol        int
	blinkDeletedRow   bool
	deletedRow        int
	currentQuery      db.Query
	shouldRerunQuery  bool
	editedQuery       string
	lastExecutedQuery string
	cellWidth         int
	columnWidths      []int // Dynamic width for each column
	detailViewMode    bool
	detailViewContent string
	detailViewScroll  int
	isTablesList      bool
	onTableSelect     func(string) tea.Cmd
	selectedTableName string
	sortColumn        string
	sortDirection     string // "", "ASC", "DESC"
}

type blinkMsg struct{}

func New(
	columns []string,
	data [][]string,
	elapsed time.Duration,
	conn db.DatabaseConnection,
	tableName, primaryKeyCol string,
	query db.Query,
	columnWidth int,
) Model {
	columnTypes := make([]string, len(columns))
	if tableName != "" && conn != nil {
		metadata, err := conn.GetTableMetadata(tableName)

		if err == nil && metadata != nil {
			colTypeMap := map[string]string{}
			for i, colName := range metadata.Columns {
				if i < len(metadata.ColumnTypes) {
					colTypeMap[colName] = metadata.ColumnTypes[i]
				}
			}
			for i, col := range columns {
				if t, ok := colTypeMap[col]; ok {
					columnTypes[i] = t
				}
			}
		}
	}

	// Extract sort information from query if present
	sortCol, sortDir := extractSortFromQuery(query.SQL)

	m := Model{
		selectedRow:      0,
		selectedCol:      0,
		offsetX:          0,
		offsetY:          0,
		columns:          columns,
		columnTypes:      columnTypes,
		data:             data,
		elapsed:          elapsed,
		visualMode:       false,
		dbConnection:     conn,
		tableName:        tableName,
		primaryKeyCol:    primaryKeyCol,
		currentQuery:     query,
		shouldRerunQuery: false,
		editedQuery:      "",
		cellWidth:        columnWidth,
		isTablesList:     false,
		sortColumn:       sortCol,
		sortDirection:    sortDir,
	}

	// Initialize column widths (will be recalculated on first resize)
	m.columnWidths = make([]int, len(columns))

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) numRows() int {
	return len(m.data)
}

func (m Model) numCols() int {
	return len(m.columns)
}

func (m Model) ShouldRerunQuery() bool {
	return m.shouldRerunQuery
}

func (m Model) GetEditedQuery() db.Query {
	updatedQuery := m.currentQuery
	if m.editedQuery != "" {
		updatedQuery.SQL = m.editedQuery
	}
	return updatedQuery
}

func (m Model) calculateHeaderLines() int {
	titleLines := 1

	var queryToDisplay string
	if m.lastExecutedQuery != "" {
		queryToDisplay = m.lastExecutedQuery
	} else {
		queryToDisplay = m.currentQuery.SQL
	}

	formattedSQL := parser.FormatSQLWithLineBreaks(queryToDisplay)
	sqlLines := strings.Count(formattedSQL, "\n") + 1

	return titleLines + sqlLines + 1
}

func (m Model) SetTablesList(onSelect func(string) tea.Cmd) Model {
	m.isTablesList = true
	m.onTableSelect = onSelect
	return m
}

func (m Model) GetSelectedTableName() string {
	return m.selectedTableName
}

// toggleSort cycles through sort states for the current column:
// For tables list: no direction → ASC → DESC → no direction
// For regular tables: no sort → ASC → DESC → no sort
func (m Model) toggleSort() (Model, tea.Cmd) {
	if m.selectedCol < 0 || m.selectedCol >= len(m.columns) {
		return m, nil
	}

	currentCol := m.columns[m.selectedCol]

	// Cycle through sort states
	if m.sortColumn != currentCol {
		// New column selected - start with ASC
		m.sortColumn = currentCol
		m.sortDirection = "ASC"
	} else if m.sortDirection == "" {
		// Same column, no direction - change to ASC
		m.sortDirection = "ASC"
	} else if m.sortDirection == "ASC" {
		// Same column, was ASC - change to DESC
		m.sortDirection = "DESC"
	} else if m.sortDirection == "DESC" {
		// Same column, was DESC
		if m.isTablesList {
			// For tables list, keep column but remove direction
			m.sortDirection = ""
		} else {
			// For regular tables, remove both column and direction
			m.sortColumn = ""
			m.sortDirection = ""
		}
	}

	// Modify the query with ORDER BY clause
	m.editedQuery = m.applySortToQuery()
	m.shouldRerunQuery = true

	return m, tea.Quit
}

// applySortToQuery adds or modifies the ORDER BY clause in the SQL query
func (m Model) applySortToQuery() string {
	sql := m.currentQuery.SQL
	if m.lastExecutedQuery != "" {
		// Use the last executed query as base if available
		sql = m.lastExecutedQuery
	}

	sql = strings.TrimSpace(sql)
	sql = strings.TrimRight(sql, ";")

	// Extract LIMIT and OFFSET clauses if present
	var limitClause string
	limitRegex := regexp.MustCompile(
		`(?i)\s+(LIMIT\s+\d+(?:\s+OFFSET\s+\d+)?)\s*$`,
	)
	if match := limitRegex.FindStringSubmatch(sql); match != nil {
		limitClause = " " + match[1]
		sql = limitRegex.ReplaceAllString(sql, "")
	}

	// Remove only the last (outermost) ORDER BY clause
	// This prevents breaking subqueries that have their own ORDER BY
	sql = removeLastOrderBy(sql)

	sql = strings.TrimSpace(sql)

	// Add new ORDER BY if we have a sort column
	if m.sortColumn != "" {
		if m.sortDirection != "" {
			sql = fmt.Sprintf(
				"%s ORDER BY %s %s",
				sql,
				m.sortColumn,
				m.sortDirection,
			)
		} else {
			// No direction means default (ASC)
			sql = fmt.Sprintf(
				"%s ORDER BY %s",
				sql,
				m.sortColumn,
			)
		}
	}

	// Re-add LIMIT clause if it existed
	if limitClause != "" {
		sql = sql + limitClause
	}

	return sql
}

// removeLastOrderBy removes only the last ORDER BY clause in a SQL query
// This is important for queries with subqueries that have their own ORDER BY
func removeLastOrderBy(sql string) string {
	// Find all occurrences of ORDER BY (case-insensitive)
	re := regexp.MustCompile(`(?i)\s+ORDER\s+BY\s+`)
	matches := re.FindAllStringIndex(sql, -1)

	if len(matches) == 0 {
		return sql
	}

	// Find the last ORDER BY that is NOT inside parentheses (i.e., not in a subquery)
	// We check by counting parentheses from the start of the string
	var lastOuterOrderByStart int = -1

	for _, match := range matches {
		orderByPos := match[0]

		// Count open and close parentheses before this ORDER BY
		openParens := strings.Count(sql[:orderByPos], "(")
		closeParens := strings.Count(sql[:orderByPos], ")")

		// If we're at the same level (not inside parentheses), this is an outer ORDER BY
		if openParens == closeParens {
			lastOuterOrderByStart = orderByPos
		}
	}

	// If no outer ORDER BY found, don't remove anything
	if lastOuterOrderByStart == -1 {
		return sql
	}

	// Find the end of the ORDER BY clause
	// It ends at: LIMIT, OFFSET, semicolon, or end of string
	remainder := sql[lastOuterOrderByStart:]
	endMarkers := regexp.MustCompile(`(?i)\s+(LIMIT|OFFSET|;)`)
	endMatch := endMarkers.FindStringIndex(remainder)

	var lastOrderByEnd int
	if endMatch != nil {
		lastOrderByEnd = lastOuterOrderByStart + endMatch[0]
	} else {
		lastOrderByEnd = len(sql)
	}

	// Remove the last outer ORDER BY clause
	return sql[:lastOuterOrderByStart] + sql[lastOrderByEnd:]
}

// extractSortFromQuery analyzes a SQL query and extracts ORDER BY information
func extractSortFromQuery(sql string) (column string, direction string) {
	// Look for ORDER BY clause with optional direction (case-insensitive)
	// We need to find the LAST (outermost) ORDER BY, not the first one
	orderByRegex := regexp.MustCompile(
		`(?i)\s+ORDER\s+BY\s+(\w+)(?:\s+(ASC|DESC))?`,
	)

	// Find all matches
	allMatches := orderByRegex.FindAllStringSubmatch(sql, -1)

	if len(allMatches) == 0 {
		return "", ""
	}

	// We want to find the last ORDER BY that is NOT inside parentheses (i.e., not in a subquery)
	// Simple heuristic: find the last match that appears after all closing parentheses

	// Find all match positions
	allMatchIndexes := orderByRegex.FindAllStringIndex(sql, -1)

	// Find the position of the last closing parenthesis
	lastParenPos := strings.LastIndex(sql, ")")

	// Find the last ORDER BY that appears after the last closing paren
	var lastOuterMatch []string
	for i := len(allMatchIndexes) - 1; i >= 0; i-- {
		matchPos := allMatchIndexes[i][0]
		if lastParenPos == -1 || matchPos > lastParenPos {
			lastOuterMatch = allMatches[i]
			break
		}
	}

	// If we didn't find one after parens, use the last match overall
	if lastOuterMatch == nil {
		lastOuterMatch = allMatches[len(allMatches)-1]
	}

	if len(lastOuterMatch) >= 2 {
		column = lastOuterMatch[1]
		if len(lastOuterMatch) >= 3 && lastOuterMatch[2] != "" {
			direction = strings.ToUpper(lastOuterMatch[2])
		}
		return column, direction
	}

	return "", ""
}

// calculateColumnWidths computes dynamic widths for each column based on content and available space
func (m *Model) calculateColumnWidths() {
	if len(m.columns) == 0 {
		return
	}

	const minWidth = 8    // Minimum column width
	const maxWidth = 50   // Maximum column width for a single column
	const borderWidth = 1 // Space for border between columns

	// Calculate content-based widths
	contentWidths := make([]int, len(m.columns))

	for i, col := range m.columns {
		// Start with header width (including icons)
		typeIcon := ""
		if i < len(m.columnTypes) && m.columnTypes[i] != "" {
			typeIcon = "Α " // 2 chars for icon
		}

		pkIcon := ""
		if m.primaryKeyCol != "" && m.columns[i] == m.primaryKeyCol {
			pkIcon = "⚿ " // 2 chars for pk icon
		}

		sortIcon := ""
		if m.sortColumn != "" && m.columns[i] == m.sortColumn {
			sortIcon = " ↑" // 2 chars for sort icon
		}

		headerLen := len([]rune(pkIcon + typeIcon + col + sortIcon))
		contentWidths[i] = headerLen

		// Sample up to 100 rows to determine content width
		sampleSize := 100
		if len(m.data) < sampleSize {
			sampleSize = len(m.data)
		}

		for j := 0; j < sampleSize; j++ {
			if i < len(m.data[j]) {
				cellLen := len([]rune(m.data[j][i]))
				if cellLen > contentWidths[i] {
					contentWidths[i] = cellLen
				}
			}
		}

		// Apply min/max constraints
		if contentWidths[i] < minWidth {
			contentWidths[i] = minWidth
		}
		if contentWidths[i] > maxWidth {
			contentWidths[i] = maxWidth
		}
	}

	// Calculate total width needed
	totalContentWidth := 0
	for _, w := range contentWidths {
		totalContentWidth += w
	}

	// Add border widths (n-1 borders for n columns)
	totalBordersWidth := (len(m.columns) - 1) * borderWidth
	totalNeededWidth := totalContentWidth + totalBordersWidth

	// Available width for table
	availableWidth := m.width - 4 // Leave some margin
	if availableWidth < 20 {
		availableWidth = 20
	}

	// Adjust widths to fit available space
	if totalNeededWidth > availableWidth {
		// Need to shrink columns proportionally
		scale := float64(
			availableWidth-totalBordersWidth,
		) / float64(
			totalContentWidth,
		)

		// First pass: scale down while respecting minimum
		scaledTotal := 0
		for i := range contentWidths {
			scaled := int(float64(contentWidths[i]) * scale)
			if scaled < minWidth {
				scaled = minWidth
			}
			contentWidths[i] = scaled
			scaledTotal += scaled
		}

		// Second pass: if still too wide, trim larger columns
		if scaledTotal+totalBordersWidth > availableWidth {
			excess := (scaledTotal + totalBordersWidth) - availableWidth

			// Sort columns by width (largest first) and trim
			for excess > 0 {
				largestIdx := 0
				largestWidth := contentWidths[0]

				for i := 1; i < len(contentWidths); i++ {
					if contentWidths[i] > largestWidth {
						largestWidth = contentWidths[i]
						largestIdx = i
					}
				}

				if contentWidths[largestIdx] > minWidth {
					contentWidths[largestIdx]--
					excess--
				} else {
					break // Can't shrink anymore
				}
			}
		}
	} else if totalNeededWidth < availableWidth {
		// Have extra space - distribute among columns
		extraSpace := availableWidth - totalNeededWidth
		perColumn := extraSpace / len(contentWidths)
		remainder := extraSpace % len(contentWidths)

		for i := range contentWidths {
			// Don't exceed max width when expanding
			addition := perColumn
			if i < remainder {
				addition++
			}

			if contentWidths[i]+addition <= maxWidth {
				contentWidths[i] += addition
			}
		}
	}

	m.columnWidths = contentWidths

	// Recalculate visible columns based on new widths
	m.calculateVisibleColumns()
}

// calculateVisibleColumns determines how many columns fit in the current view
func (m *Model) calculateVisibleColumns() {
	if len(m.columnWidths) == 0 {
		m.visibleCols = 0
		return
	}

	availableWidth := m.width - 4
	usedWidth := 0
	visibleCount := 0

	for i := m.offsetX; i < len(m.columnWidths); i++ {
		colWidth := m.columnWidths[i]
		borderWidth := 0
		if i > m.offsetX {
			borderWidth = 1
		}

		if usedWidth+colWidth+borderWidth <= availableWidth {
			usedWidth += colWidth + borderWidth
			visibleCount++
		} else {
			break
		}
	}

	if visibleCount == 0 && len(m.columnWidths) > 0 {
		visibleCount = 1 // Always show at least one column
	}

	m.visibleCols = visibleCount
}

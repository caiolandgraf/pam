package table

import (
	"testing"
	"time"

	"github.com/caiolandgraf/pam/internal/db"
)

func TestNew(t *testing.T) {
	columns := []string{"id", "name", "email"}
	data := [][]string{
		{"1", "Alice", "alice@example.com"},
		{"2", "Bob", "bob@example.com"},
	}
	elapsed := 100 * time.Millisecond
	query := db.Query{
		Name: "test_query",
		SQL:  "SELECT * FROM users",
		Id:   1,
	}

	model := New(columns, nil, data, elapsed, nil, "", "", query, 15)

	// Verify initial state
	if model.selectedRow != 0 {
		t.Errorf("New() selectedRow = %d, want 0", model.selectedRow)
	}
	if model.selectedCol != 0 {
		t.Errorf("New() selectedCol = %d, want 0", model.selectedCol)
	}
	if len(model.columns) != len(columns) {
		t.Errorf(
			"New() columns length = %d, want %d",
			len(model.columns),
			len(columns),
		)
	}
	if len(model.data) != len(data) {
		t.Errorf("New() data length = %d, want %d", len(model.data), len(data))
	}
	if model.currentQuery.Name != "test_query" {
		t.Errorf(
			"New() query name = %s, want test_query",
			model.currentQuery.Name,
		)
	}
	if model.cellWidth != 15 {
		t.Errorf("New() cellWidth = %d, want 15", model.cellWidth)
	}
	if model.visualMode {
		t.Error("New() visualMode should be false initially")
	}
	if model.detailViewMode {
		t.Error("New() detailViewMode should be false initially")
	}
}

func TestExtractSortFromQuery(t *testing.T) {
	tests := []struct {
		name          string
		sql           string
		wantColumn    string
		wantDirection string
	}{
		{
			name:          "ORDER BY ASC",
			sql:           "SELECT * FROM users ORDER BY name ASC",
			wantColumn:    "name",
			wantDirection: "ASC",
		},
		{
			name:          "ORDER BY DESC",
			sql:           "SELECT * FROM users ORDER BY created_at DESC",
			wantColumn:    "created_at",
			wantDirection: "DESC",
		},
		{
			name:          "ORDER BY without direction (defaults to ASC)",
			sql:           "SELECT * FROM users ORDER BY id",
			wantColumn:    "id",
			wantDirection: "ASC",
		},
		{
			name:          "No ORDER BY clause",
			sql:           "SELECT * FROM users",
			wantColumn:    "",
			wantDirection: "",
		},
		{
			name:          "ORDER BY with table prefix",
			sql:           "SELECT * FROM users u ORDER BY u.name ASC",
			wantColumn:    "name",
			wantDirection: "ASC",
		},
		{
			name:          "Multiple ORDER BY (takes first)",
			sql:           "SELECT * FROM users ORDER BY name ASC, id DESC",
			wantColumn:    "name",
			wantDirection: "ASC",
		},
		{
			name:          "ORDER BY with mixed case",
			sql:           "SELECT * FROM users order by Name DESC",
			wantColumn:    "name",
			wantDirection: "DESC",
		},
		{
			name:          "Complex query with ORDER BY",
			sql:           "SELECT u.id, u.name FROM users u WHERE u.active = true ORDER BY u.name DESC LIMIT 10",
			wantColumn:    "name",
			wantDirection: "DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col, dir := extractSortFromQuery(tt.sql)
			if col != tt.wantColumn {
				t.Errorf(
					"extractSortFromQuery() column = %s, want %s",
					col,
					tt.wantColumn,
				)
			}
			if dir != tt.wantDirection {
				t.Errorf(
					"extractSortFromQuery() direction = %s, want %s",
					dir,
					tt.wantDirection,
				)
			}
		})
	}
}

func TestModel_GetEditedQuery(t *testing.T) {
	query := db.Query{
		Name: "test_query",
		SQL:  "SELECT * FROM users",
		Id:   1,
	}
	model := New(
		[]string{"id"},
		nil,
		[][]string{{"1"}},
		0,
		nil,
		"",
		"",
		query,
		15,
	)

	// Initially should return current query
	editedQuery := model.GetEditedQuery()
	if editedQuery.SQL != query.SQL {
		t.Errorf(
			"GetEditedQuery() SQL = %s, want %s",
			editedQuery.SQL,
			query.SQL,
		)
	}

	// After editing
	model.editedQuery = "SELECT * FROM users WHERE active = true"
	editedQuery = model.GetEditedQuery()
	if editedQuery.SQL != model.editedQuery {
		t.Errorf(
			"GetEditedQuery() SQL = %s, want %s",
			editedQuery.SQL,
			model.editedQuery,
		)
	}
}

func TestModel_ShouldRerunQuery(t *testing.T) {
	model := New(
		[]string{"id"},
		nil,
		[][]string{{"1"}},
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	// Initially should not rerun
	if model.ShouldRerunQuery() {
		t.Error("ShouldRerunQuery() should be false initially")
	}

	// After setting flag
	model.shouldRerunQuery = true
	if !model.ShouldRerunQuery() {
		t.Error("ShouldRerunQuery() should be true after setting flag")
	}
}

func TestModel_GetSelectedTableName(t *testing.T) {
	model := New(
		[]string{"name"},
		nil,
		[][]string{{"users"}},
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	// Initially empty
	if model.GetSelectedTableName() != "" {
		t.Errorf(
			"GetSelectedTableName() = %s, want empty string",
			model.GetSelectedTableName(),
		)
	}

	// After selecting
	model.selectedTableName = "users"
	if model.GetSelectedTableName() != "users" {
		t.Errorf(
			"GetSelectedTableName() = %s, want users",
			model.GetSelectedTableName(),
		)
	}
}

func TestModel_ExtractNewValue(t *testing.T) {
	model := New(
		[]string{"id", "name", "email"},
		nil,
		[][]string{},
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	tests := []struct {
		name       string
		sql        string
		columnName string
		want       string
	}{
		{
			name:       "single quote value",
			sql:        "UPDATE users SET name = 'John Doe' WHERE id = 1",
			columnName: "name",
			want:       "John Doe",
		},
		{
			name:       "numeric value",
			sql:        "UPDATE users SET age = 30 WHERE id = 1",
			columnName: "age",
			want:       "30",
		},
		{
			name:       "value with spaces",
			sql:        "UPDATE users SET description = 'This is a test' WHERE id = 1",
			columnName: "description",
			want:       "This is a test",
		},
		{
			name:       "value with special characters",
			sql:        "UPDATE users SET email = 'test@example.com' WHERE id = 1",
			columnName: "email",
			want:       "test@example.com",
		},
		{
			name:       "NULL value",
			sql:        "UPDATE users SET middle_name = NULL WHERE id = 1",
			columnName: "middle_name",
			want:       "NULL",
		},
		{
			name:       "boolean value",
			sql:        "UPDATE users SET active = true WHERE id = 1",
			columnName: "active",
			want:       "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.extractNewValue(tt.sql, tt.columnName)
			if result != tt.want {
				t.Errorf("extractNewValue() = %s, want %s", result, tt.want)
			}
		})
	}
}

func TestModel_BlinkCmd(t *testing.T) {
	model := New(
		[]string{"id"},
		nil,
		[][]string{{"1"}},
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	cmd := model.blinkCmd()
	if cmd == nil {
		t.Error("blinkCmd() should return a non-nil command")
	}

	// Execute the command and verify it returns blinkMsg
	msg := cmd()
	if _, ok := msg.(blinkMsg); !ok {
		t.Errorf("blinkCmd() message type = %T, want blinkMsg", msg)
	}
}

func TestModel_NavigationBounds(t *testing.T) {
	data := [][]string{
		{"1", "Alice", "alice@example.com"},
		{"2", "Bob", "bob@example.com"},
		{"3", "Charlie", "charlie@example.com"},
	}
	model := New(
		[]string{"id", "name", "email"},
		nil,
		data,
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	// Test row bounds
	if model.selectedRow < 0 {
		t.Error("selectedRow should not be negative")
	}
	if model.selectedRow >= len(data) {
		t.Errorf(
			"selectedRow %d should be less than data length %d",
			model.selectedRow,
			len(data),
		)
	}

	// Test column bounds
	if model.selectedCol < 0 {
		t.Error("selectedCol should not be negative")
	}
	if model.selectedCol >= len(model.columns) {
		t.Errorf(
			"selectedCol %d should be less than columns length %d",
			model.selectedCol,
			len(model.columns),
		)
	}
}

func TestModel_VisualModeToggle(t *testing.T) {
	model := New(
		[]string{"id"},
		nil,
		[][]string{{"1"}},
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	// Initially not in visual mode
	if model.visualMode {
		t.Error("visualMode should be false initially")
	}

	// Toggle visual mode
	model.visualMode = true
	model.visualStartRow = 0
	model.visualStartCol = 0

	if !model.visualMode {
		t.Error("visualMode should be true after toggling")
	}
	if model.visualStartRow != 0 {
		t.Errorf("visualStartRow = %d, want 0", model.visualStartRow)
	}
	if model.visualStartCol != 0 {
		t.Errorf("visualStartCol = %d, want 0", model.visualStartCol)
	}
}

func TestModel_DetailViewToggle(t *testing.T) {
	model := New(
		[]string{"id"},
		nil,
		[][]string{{"1"}},
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	// Initially not in detail view
	if model.detailViewMode {
		t.Error("detailViewMode should be false initially")
	}

	// Enter detail view
	model.detailViewMode = true
	model.detailViewContent = "test content"
	model.detailViewScroll = 0

	if !model.detailViewMode {
		t.Error("detailViewMode should be true")
	}
	if model.detailViewContent != "test content" {
		t.Errorf(
			"detailViewContent = %s, want 'test content'",
			model.detailViewContent,
		)
	}
	if model.detailViewScroll != 0 {
		t.Errorf("detailViewScroll = %d, want 0", model.detailViewScroll)
	}
}

func TestModel_CellWidth(t *testing.T) {
	tests := []struct {
		name      string
		cellWidth int
		wantWidth int
	}{
		{
			name:      "default cell width",
			cellWidth: 15,
			wantWidth: 15,
		},
		{
			name:      "custom cell width",
			cellWidth: 25,
			wantWidth: 25,
		},
		{
			name:      "zero cell width should be allowed",
			cellWidth: 0,
			wantWidth: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New(
				[]string{"id"},
				nil,
				[][]string{{"1"}},
				0,
				nil,
				"",
				"",
				db.Query{},
				tt.cellWidth,
			)
			if model.cellWidth != tt.wantWidth {
				t.Errorf(
					"cellWidth = %d, want %d",
					model.cellWidth,
					tt.wantWidth,
				)
			}
		})
	}
}

func TestModel_SortDirection(t *testing.T) {
	tests := []struct {
		name          string
		query         db.Query
		wantColumn    string
		wantDirection string
	}{
		{
			name: "query with ASC sort",
			query: db.Query{
				Name: "sorted_asc",
				SQL:  "SELECT * FROM users ORDER BY name ASC",
			},
			wantColumn:    "name",
			wantDirection: "ASC",
		},
		{
			name: "query with DESC sort",
			query: db.Query{
				Name: "sorted_desc",
				SQL:  "SELECT * FROM users ORDER BY created_at DESC",
			},
			wantColumn:    "created_at",
			wantDirection: "DESC",
		},
		{
			name: "query without sort",
			query: db.Query{
				Name: "unsorted",
				SQL:  "SELECT * FROM users",
			},
			wantColumn:    "",
			wantDirection: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New(
				[]string{"id", "name"},
				nil,
				[][]string{},
				0,
				nil,
				"",
				"",
				tt.query,
				15,
			)
			if model.sortColumn != tt.wantColumn {
				t.Errorf(
					"sortColumn = %s, want %s",
					model.sortColumn,
					tt.wantColumn,
				)
			}
			if model.sortDirection != tt.wantDirection {
				t.Errorf(
					"sortDirection = %s, want %s",
					model.sortDirection,
					tt.wantDirection,
				)
			}
		})
	}
}

func TestModel_EmptyData(t *testing.T) {
	// Test with empty data
	model := New(
		[]string{"id", "name"},
		nil,
		[][]string{},
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	if len(model.data) != 0 {
		t.Errorf("data length = %d, want 0", len(model.data))
	}
	if model.selectedRow != 0 {
		t.Errorf("selectedRow = %d, want 0", model.selectedRow)
	}

	// Should handle navigation gracefully without panicking
	// This is important for empty result sets
}

func TestModel_SingleRow(t *testing.T) {
	// Test with single row
	data := [][]string{{"1", "Alice"}}
	model := New(
		[]string{"id", "name"},
		nil,
		data,
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	if len(model.data) != 1 {
		t.Errorf("data length = %d, want 1", len(model.data))
	}
	if model.selectedRow != 0 {
		t.Errorf("selectedRow = %d, want 0", model.selectedRow)
	}
}

func TestModel_LargeDataset(t *testing.T) {
	// Test with large dataset
	data := make([][]string, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []string{string(rune(i)), "name", "email"}
	}

	model := New(
		[]string{"id", "name", "email"},
		nil,
		data,
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	if len(model.data) != 1000 {
		t.Errorf("data length = %d, want 1000", len(model.data))
	}
	if model.selectedRow != 0 {
		t.Errorf("selectedRow = %d, want 0", model.selectedRow)
	}
}

func TestModel_BlinkStates(t *testing.T) {
	model := New(
		[]string{"id"},
		nil,
		[][]string{{"1"}},
		0,
		nil,
		"",
		"",
		db.Query{},
		15,
	)

	// Test blink states are initially false
	if model.blinkCopiedCell {
		t.Error("blinkCopiedCell should be false initially")
	}
	if model.blinkUpdatedCell {
		t.Error("blinkUpdatedCell should be false initially")
	}
	if model.blinkDeletedRow {
		t.Error("blinkDeletedRow should be false initially")
	}

	// Test setting blink states
	model.blinkCopiedCell = true
	if !model.blinkCopiedCell {
		t.Error("blinkCopiedCell should be true after setting")
	}

	model.blinkUpdatedCell = true
	model.updatedRow = 1
	model.updatedCol = 2
	if !model.blinkUpdatedCell {
		t.Error("blinkUpdatedCell should be true after setting")
	}
	if model.updatedRow != 1 || model.updatedCol != 2 {
		t.Error("updated row/col should be set correctly")
	}

	model.blinkDeletedRow = true
	model.deletedRow = 3
	if !model.blinkDeletedRow {
		t.Error("blinkDeletedRow should be true after setting")
	}
	if model.deletedRow != 3 {
		t.Error("deletedRow should be set correctly")
	}
}

func TestToggleSort_Cycle(t *testing.T) {
	// Test the 3-state cycle: no sort → ASC → DESC → no sort
	// for both regular tables and tables list
	for _, isTablesList := range []bool{false, true} {
		t.Run(func() string {
			if isTablesList {
				return "tables_list"
			}
			return "regular_table"
		}(), func(t *testing.T) {
			model := New(
				[]string{"id", "name"},
				nil,
				[][]string{{"1", "Alice"}, {"2", "Bob"}},
				0,
				nil,
				"",
				"",
				db.Query{SQL: "SELECT id, name FROM users"},
				15,
			)
			model.isTablesList = isTablesList
			model.selectedCol = 1 // "name" column

			// State 0: no sort
			if model.sortColumn != "" || model.sortDirection != "" {
				t.Fatalf("initial state: want no sort, got col=%q dir=%q",
					model.sortColumn, model.sortDirection)
			}

			// F press 1: no sort → ASC on "name"
			model, _ = model.toggleSort()
			if model.sortColumn != "name" || model.sortDirection != "ASC" {
				t.Errorf(
					"after 1st F: want col=name dir=ASC, got col=%q dir=%q",
					model.sortColumn,
					model.sortDirection,
				)
			}
			if model.editedQuery == "" {
				t.Error("after 1st F: editedQuery should not be empty")
			}
			if !model.shouldRerunQuery {
				t.Error("after 1st F: shouldRerunQuery should be true")
			}

			// Simulate re-run: new model with sorted SQL, same col selected
			model = New(
				[]string{"id", "name"},
				nil,
				[][]string{{"1", "Alice"}, {"2", "Bob"}},
				0,
				nil,
				"",
				"",
				db.Query{SQL: model.editedQuery},
				15,
			)
			model.isTablesList = isTablesList
			model.selectedCol = 1

			// F press 2: ASC → DESC
			model, _ = model.toggleSort()
			if model.sortColumn != "name" || model.sortDirection != "DESC" {
				t.Errorf(
					"after 2nd F: want col=name dir=DESC, got col=%q dir=%q",
					model.sortColumn,
					model.sortDirection,
				)
			}

			// Simulate re-run again
			model = New(
				[]string{"id", "name"},
				nil,
				[][]string{{"1", "Alice"}, {"2", "Bob"}},
				0,
				nil,
				"",
				"",
				db.Query{SQL: model.editedQuery},
				15,
			)
			model.isTablesList = isTablesList
			model.selectedCol = 1

			// F press 3: DESC → no sort (clear both)
			model, _ = model.toggleSort()
			if model.sortColumn != "" || model.sortDirection != "" {
				t.Errorf("after 3rd F: want no sort, got col=%q dir=%q",
					model.sortColumn, model.sortDirection)
			}

			// Simulate re-run again
			model = New(
				[]string{"id", "name"},
				nil,
				[][]string{{"1", "Alice"}, {"2", "Bob"}},
				0,
				nil,
				"",
				"",
				db.Query{SQL: model.editedQuery},
				15,
			)
			model.isTablesList = isTablesList
			model.selectedCol = 1

			// After re-run with no ORDER BY, sort state should be empty
			if model.sortColumn != "" || model.sortDirection != "" {
				t.Errorf("after 3rd F re-run: want no sort, got col=%q dir=%q",
					model.sortColumn, model.sortDirection)
			}

			// F press 4: back to ASC (cycle restarts cleanly)
			model, _ = model.toggleSort()
			if model.sortColumn != "name" || model.sortDirection != "ASC" {
				t.Errorf(
					"after 4th F: want col=name dir=ASC, got col=%q dir=%q",
					model.sortColumn,
					model.sortDirection,
				)
			}
		})
	}
}

func TestToggleSort_TablesListInitialASC(t *testing.T) {
	// Simulates the tables list scenario: initial query already has ORDER BY name ASC
	// (built into nameOnlyQuery). Pressing F should cycle DESC → no sort → ASC, not loop.
	initialSQL := "SELECT name FROM (SELECT TABLE_NAME as name FROM information_schema.TABLES) AS t ORDER BY name ASC"

	model := New(
		[]string{"name"},
		nil,
		[][]string{{"users"}, {"orders"}},
		0,
		nil,
		"",
		"",
		db.Query{SQL: initialSQL},
		15,
	)
	model.isTablesList = true
	model.selectedCol = 0 // "name"

	// Initial state parsed from query: should be ASC
	if model.sortColumn != "name" || model.sortDirection != "ASC" {
		t.Fatalf("initial parse: want col=name dir=ASC, got col=%q dir=%q",
			model.sortColumn, model.sortDirection)
	}

	// F press 1: ASC → DESC (not a loop back to ASC)
	model, _ = model.toggleSort()
	if model.sortDirection != "DESC" {
		t.Errorf("after 1st F: want DESC, got %q", model.sortDirection)
	}
	if !contains(model.editedQuery, "ORDER BY name DESC") {
		t.Errorf(
			"after 1st F: editedQuery should contain ORDER BY name DESC, got: %s",
			model.editedQuery,
		)
	}

	// Simulate re-run with DESC query
	model = New(
		[]string{"name"},
		nil,
		[][]string{{"users"}, {"orders"}},
		0,
		nil,
		"",
		"",
		db.Query{SQL: model.editedQuery},
		15,
	)
	model.isTablesList = true
	model.selectedCol = 0

	// F press 2: DESC → no sort (not back to ASC — this was the bug)
	model, _ = model.toggleSort()
	if model.sortColumn != "" || model.sortDirection != "" {
		t.Errorf(
			"after 2nd F: want no sort (clear), got col=%q dir=%q — infinite loop bug?",
			model.sortColumn,
			model.sortDirection,
		)
	}
	// The SQL should NOT contain an ORDER BY at the outer level
	if contains(model.editedQuery, "ORDER BY name") {
		t.Errorf(
			"after 2nd F: editedQuery should not contain outer ORDER BY name, got: %s",
			model.editedQuery,
		)
	}
}

func TestToggleSort_DifferentColumn(t *testing.T) {
	// Switching to a different column should restart the cycle at ASC
	model := New(
		[]string{"id", "name", "email"},
		nil,
		[][]string{{"1", "Alice", "a@b.com"}},
		0,
		nil,
		"",
		"",
		db.Query{SQL: "SELECT id, name, email FROM users ORDER BY name ASC"},
		15,
	)
	model.selectedCol = 1 // "name" — already sorted

	if model.sortColumn != "name" || model.sortDirection != "ASC" {
		t.Fatalf("setup: want col=name dir=ASC, got col=%q dir=%q",
			model.sortColumn, model.sortDirection)
	}

	// Move to a different column and press F
	model.selectedCol = 2 // "email"
	model, _ = model.toggleSort()

	if model.sortColumn != "email" || model.sortDirection != "ASC" {
		t.Errorf("switching column: want col=email dir=ASC, got col=%q dir=%q",
			model.sortColumn, model.sortDirection)
	}
	if !contains(model.editedQuery, "ORDER BY email ASC") {
		t.Errorf(
			"switching column: editedQuery should contain ORDER BY email ASC, got: %s",
			model.editedQuery,
		)
	}
}

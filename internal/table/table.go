package table

import (
	"time"

	"github.com/caiolandgraf/pam/internal/config"
	"github.com/caiolandgraf/pam/internal/db"
	tea "github.com/charmbracelet/bubbletea"
)

func Render(
	columns []string,
	columnTypes []string,
	data [][]string,
	elapsed time.Duration,
	conn db.DatabaseConnection,
	tableName, primaryKeyCol string,
	query db.Query,
	columnWidth int,
	visibility config.UIVisibility,
	saveCallback func(query db.Query) (db.Query, error),
	initialStatus ...string,
) (Model, error) {
	model := New(
		columns,
		columnTypes,
		data,
		elapsed,
		conn,
		tableName,
		primaryKeyCol,
		query,
		columnWidth,
		visibility,
	)
	model.saveQueryCallback = saveCallback
	if len(initialStatus) > 0 && initialStatus[0] != "" {
		model.statusMessage = initialStatus[0]
	}
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return model, err
	}
	return finalModel.(Model), nil
}

func RenderTablesList(
	columns []string,
	data [][]string,
	elapsed time.Duration,
	conn db.DatabaseConnection,
	query db.Query,
	columnWidth int,
	visibility config.UIVisibility,
) (Model, error) {
	model := New(columns, nil, data, elapsed, conn, "", "", query, columnWidth, visibility)
	model.isTablesList = true
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return model, err
	}
	return finalModel.(Model), nil
}

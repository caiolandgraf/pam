package table

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/eduardofuncao/pam/internal/db"
)

func Render(columns []string, data [][]string, elapsed time.Duration, conn db. DatabaseConnection, tableName, primaryKeyCol string, query db.Query) (Model, error) {
	model := New(columns, data, elapsed, conn, tableName, primaryKeyCol, query)
	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return model, err
	}
	return finalModel.(Model), nil
}

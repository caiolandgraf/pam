package table

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) deleteRow() (tea.Model, tea.Cmd) {
	if m.selectedRow < 0 || m. selectedRow >= m.numRows() {
		return m, nil
	}

	if m.primaryKeyCol == "" {
		fmt.Fprintf(os.Stderr, "\nCannot delete:  No primary key defined for this table\n")
		tea.EnterAltScreen()
		return m, nil
	}

	deleteStmt := m.buildDeleteStatement()

	editorCmd := os.Getenv("EDITOR")
	if editorCmd == "" {
		editorCmd = "vim"
	}

	tmpFile, err := os.CreateTemp("", "pam-delete-*.sql")
	if err != nil {
		fmt. Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
		return m, tea.Quit
	}
	tmpPath := tmpFile. Name()
	defer os.Remove(tmpPath)

	header := `-- DELETE Statement
-- WARNING: This will permanently delete data! 
-- Ensure the WHERE clause is present and correct before saving. 
--
`
	content := header + deleteStmt

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		tmpFile.Close()
		fmt.Fprintf(os. Stderr, "Error writing temp file: %v\n", err)
		return m, tea. Quit
	}
	tmpFile.Close()

	tea.ExitAltScreen()

	cmd := exec.Command(editorCmd, tmpPath)
	cmd.Stdin = os. Stdin
	cmd.Stdout = os.Stdout
	cmd. Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt. Fprintf(os.Stderr, "Error running editor: %v\n", err)
		tea.EnterAltScreen()
		return m, tea.Quit
	}

	editedSQL, err := os.ReadFile(tmpPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading edited file: %v\n", err)
		tea.EnterAltScreen()
		return m, tea.Quit
	}

	tea.EnterAltScreen()

	sqlStr := string(editedSQL)

	if err := validateDeleteStatement(sqlStr); err != nil {
		fmt. Fprintf(os.Stderr, "\nDelete validation failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "Press any key to continue.. .\n")
		var buf [1]byte
		os.Stdin.Read(buf[: ])
		return m, nil
	}

	if err := m.executeDelete(sqlStr); err != nil {
		fmt.Fprintf(os.Stderr, "\nError executing delete: %v\n", err)
		fmt.Fprintf(os.Stderr, "Press any key to continue.. .\n")
		var buf [1]byte
		os.Stdin.Read(buf[: ])
		return m, nil
	}

	m.data = append(m.data[:m.selectedRow], m.data[m.selectedRow+1:]...)

	if m.selectedRow >= m.numRows() && m.numRows() > 0 {
		m.selectedRow = m.numRows() - 1
	}

	if m.offsetY >= m.numRows() && m.numRows() > 0 {
		m.offsetY = m.numRows() - 1
	}

	m.blinkDeletedRow = true
	m.deletedRow = m.selectedRow

	return m, m.blinkCmd()
}

func (m Model) buildDeleteStatement() string {
	pkValue := ""
	if m.primaryKeyCol != "" {
		for i, col := range m.columns {
			if col == m.primaryKeyCol {
				pkValue = m.data[m.selectedRow][i]
				break
			}
		}
	}

	return m.dbConnection.BuildDeleteStatement(
		m.tableName,
		m.primaryKeyCol,
		pkValue,
	)
}

func (m Model) executeDelete(sql string) error {
	var result strings.Builder
	for line := range strings.SplitSeq(sql, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "--") && trimmed != "" {
			result. WriteString(trimmed)
			result.WriteString(" ")
		}
	}

	cleanSQL := strings.TrimSpace(result.String())
	cleanSQL = strings.TrimSuffix(cleanSQL, ";")

	if cleanSQL == "" {
		return fmt.Errorf("no SQL to execute")
	}

	return m.dbConnection. Exec(cleanSQL)
}

func validateDeleteStatement(sql string) error {
	var result strings. Builder
	for line := range strings.SplitSeq(sql, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "--") && trimmed != "" {
			result.WriteString(trimmed)
			result.WriteString(" ")
		}
	}
	cleanSQL := strings.TrimSpace(result.String())

	if cleanSQL == "" {
		return fmt.Errorf("empty SQL statement")
	}

	deleteRegex := regexp.MustCompile(`(?i)^\s*DELETE\s+FROM`)
	if ! deleteRegex.MatchString(cleanSQL) {
		return fmt.Errorf("not a DELETE statement")
	}

	whereRegex := regexp.MustCompile(`(?i)\bWHERE\b`)
	if !whereRegex.MatchString(cleanSQL) {
		return fmt.Errorf("DELETE statement must include a WHERE clause for safety")
	}

	return nil
}

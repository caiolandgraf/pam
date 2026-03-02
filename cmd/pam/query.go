package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/caiolandgraf/pam/internal/config"
	"github.com/caiolandgraf/pam/internal/db"
	"github.com/caiolandgraf/pam/internal/editor"
	"github.com/caiolandgraf/pam/internal/run"
)

type queryFlags struct {
	tableName string
	editMode  bool
	sql       string
}

func parseQueryFlags() queryFlags {
	flags := queryFlags{}
	args := os.Args[2:]

	var positionals []string

	i := 0
	for i < len(args) {
		arg := args[i]

		switch {
		case arg == "--table" || arg == "-t":
			if i+1 < len(args) {
				flags.tableName = args[i+1]
				i += 2
			} else {
				printError("--table requires a value")
				i++
			}
		case strings.HasPrefix(arg, "--table="):
			flags.tableName = strings.TrimPrefix(arg, "--table=")
			i++
		case strings.HasPrefix(arg, "-t="):
			flags.tableName = strings.TrimPrefix(arg, "-t=")
			i++
		case arg == "--edit" || arg == "-e":
			flags.editMode = true
			i++
		default:
			if !strings.HasPrefix(arg, "-") {
				positionals = append(positionals, arg)
			}
			i++
		}
	}

	if len(positionals) > 0 {
		flags.sql = strings.Join(positionals, " ")
	}

	return flags
}

func (a *App) handleQuery() {
	if a.config.CurrentConnection == "" {
		printError(
			"No active connection. Use 'pam switch <connection>' or 'pam init' first",
		)
	}

	flags := parseQueryFlags()

	if flags.tableName == "" && flags.sql == "" {
		// No table and no SQL — delegate to handleRun for backward compat
		a.handleRun()
		return
	}

	conn := config.FromConnectionYaml(
		a.config.Connections[a.config.CurrentConnection],
	)

	// Build the SQL
	sql := flags.sql
	if sql == "" && flags.tableName != "" {
		// No SQL provided, default to SELECT * FROM <table>
		sql = fmt.Sprintf("SELECT * FROM %s", flags.tableName)
	}

	// If edit mode, open editor with current SQL
	if flags.editMode {
		instructions := "-- Edit your SQL query below\n-- Save and exit to execute, or clear contents to cancel\n--\n"
		edited, err := editor.EditTempFileWithTemplate(
			instructions+sql,
			"pam-query-",
		)
		if err != nil {
			printError("Error opening editor: %v", err)
		}
		if edited == "" {
			printError("Empty SQL, cancelled")
		}
		sql = edited
	}

	// Build query object
	queryName := flags.tableName
	if queryName == "" {
		queryName = "<query>"
	}

	query := db.Query{
		Name:      queryName,
		SQL:       sql,
		TableName: flags.tableName,
		Id:        -1,
	}

	if err := conn.Open(); err != nil {
		printError(
			"Could not open connection to %s/%s: %s",
			conn.GetDbType(),
			conn.GetName(),
			err,
		)
	}
	defer conn.Close()

	// Try to get primary key from table metadata if table is specified
	if flags.tableName != "" {
		if metadata, err := conn.GetTableMetadata(
			flags.tableName,
		); err == nil && metadata != nil {
			query.PrimaryKeys = metadata.PrimaryKeys
		}
	}

	if run.IsSelectQuery(sql) {
		var onRerun func(editedSQL string)
		onRerun = func(editedSQL string) {
			editedQuery := db.Query{
				Name:      queryName,
				SQL:       editedSQL,
				TableName: flags.tableName,
				Id:        -1,
			}
			run.ExecuteSelect(
				editedSQL,
				queryName,
				run.ExecutionParams{
					Query:        editedQuery,
					Connection:   conn,
					Config:       a.config,
					SaveCallback: a.saveQueryFromTable,
					OnRerun:      onRerun,
				},
			)
		}
		run.ExecuteSelect(
			sql,
			queryName,
			run.ExecutionParams{
				Query:        query,
				Connection:   conn,
				Config:       a.config,
				SaveCallback: a.saveQueryFromTable,
				OnRerun:      onRerun,
			},
		)
	} else {
		run.ExecuteNonSelect(run.ExecutionParams{
			Query:      query,
			Connection: conn,
			Config:     a.config,
		})
	}
}

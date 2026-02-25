package main

import (
	"fmt"
	"os"
	"time"

	"github.com/eduardofuncao/pam/internal/config"
	"github.com/eduardofuncao/pam/internal/spinner"
	"github.com/eduardofuncao/pam/internal/table"
)

func (a *App) handleTableView() {
	if a.config.CurrentConnection == "" {
		printError(
			"No active connection. Use 'pam switch <connection>' or 'pam init' first",
		)
	}

	args := os.Args[2:]
	if len(args) == 0 {
		printError("Usage: pam table-view <table-name>")
	}

	tableName := args[0]

	conn := config.FromConnectionYaml(
		a.config.Connections[a.config.CurrentConnection],
	)

	if err := conn.Open(); err != nil {
		printError(
			"Could not open connection to %s/%s: %s",
			conn.GetDbType(),
			conn.GetName(),
			err,
		)
	}
	defer conn.Close()

	start := time.Now()
	done := make(chan struct{})
	go spinner.CircleWaitWithTimer(done)

	columns, err := conn.GetColumnDetails(tableName)
	if err != nil {
		done <- struct{}{}
		printError("Could not get column details for '%s': %v", tableName, err)
	}

	done <- struct{}{}
	elapsed := time.Since(start)

	if len(columns) == 0 {
		fmt.Printf("No columns found for table '%s'\n", tableName)
		return
	}

	_, err = table.RenderTableView(tableName, columns, conn, elapsed)
	if err != nil {
		printError("Error rendering table view: %v", err)
	}
}

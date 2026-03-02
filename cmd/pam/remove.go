package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/caiolandgraf/pam/internal/db"
	"github.com/caiolandgraf/pam/internal/styles"
)

func (a *App) handleRemove() {
	if len(os.Args) < 3 {
		printError(
			"Usage:  pam remove <run-name>  |  pam remove --conn <connection-name>",
		)
	}

	args := os.Args[2:]

	// Check if removing a connection via --conn / -c
	connName := ""
	for i, arg := range args {
		if (arg == "--conn" || arg == "-c") && i+1 < len(args) {
			connName = args[i+1]
			break
		}
		if strings.HasPrefix(arg, "--conn=") {
			connName = strings.TrimPrefix(arg, "--conn=")
			break
		}
	}

	if connName != "" {
		a.handleRemoveConnection(connName)
		return
	}

	// Default: remove a saved query
	conn := a.config.Connections[a.config.CurrentConnection]
	queries := conn.Queries

	query, exists := db.FindQueryWithSelector(queries, os.Args[2])
	if !exists {
		printError("Query '%s' could not be found", os.Args[2])
	}

	delete(conn.Queries, query.Name)

	err := a.config.Save()
	if err != nil {
		printError("Could not save configuration file: %v", err)
	}

	fmt.Println(
		styles.Success.Render(fmt.Sprintf("✓ Removed query '%s'", query.Name)),
	)
}

func (a *App) handleRemoveConnection(name string) {
	if _, exists := a.config.Connections[name]; !exists {
		printError("Connection '%s' could not be found", name)
	}

	delete(a.config.Connections, name)

	// If we just removed the active connection, clear it
	if a.config.CurrentConnection == name {
		a.config.CurrentConnection = ""
	}

	err := a.config.Save()
	if err != nil {
		printError("Could not save configuration file: %v", err)
	}

	fmt.Println(
		styles.Success.Render(fmt.Sprintf("✓ Removed connection '%s'", name)),
	)

	if a.config.CurrentConnection == "" && len(a.config.Connections) > 0 {
		fmt.Println(
			styles.Faint.Render(
				"  No active connection. Use 'pam switch <name>' to select one.",
			),
		)
	}
}

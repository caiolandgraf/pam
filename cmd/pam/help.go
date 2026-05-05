package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/caiolandgraf/pam/internal/styles"
)

func (a *App) handleHelp() {
	if len(os.Args) == 2 {
		a.PrintGeneralHelp()
	} else {
		a.PrintCommandHelp()
	}
}

// asciiLogo returns the PAM ASCII art styled in pink
func asciiLogo() string {
	pink := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)

	lines := []string{
		`  █▀█ █▀█ █▀▄▀█`,
		`  █▀▀ █▀█ █░▀░█`,
	}

	var b strings.Builder
	for _, line := range lines {
		b.WriteString(pink.Render(line))
		b.WriteString("\n")
	}
	return b.String()
}

func separator() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
		Render(strings.Repeat("─", 52))
}

func cmdEntry(name, args, desc string) string {
	accent := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Accent)).
		Bold(true)
	faint := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Muted))
	normal := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Normal))

	nameCol := fmt.Sprintf("  pam %-14s", name)
	argsCol := fmt.Sprintf("%-18s", args)

	return accent.Render(nameCol) + faint.Render(argsCol) + normal.Render(desc)
}

func sectionHeader(title string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Primary)).
		Bold(true).
		Render(title)
}

func (a *App) PrintGeneralHelp() {
	// Header
	fmt.Println(
		styles.Title.Render(
			"PAM — Pam's Database Drawer",
		),
	)
	fmt.Println(
		styles.Faint.Render(
			"Save, edit, and run named SQL queries across connections.",
		),
	)
	fmt.Println()

	// Usage
	fmt.Println(styles.Title.Render("Usage"))
	fmt.Println(styles.Separator.Render("  pam <command> [arguments]"))
	fmt.Println()

	// ── CONNECTIONS ──────────────────────────────────────────────
	fmt.Println(sectionHeader("CONNECTIONS"))
	fmt.Println(
		cmdEntry(
			"init",
			"",
			"Create a new database connection (interactive TUI if no args)",
		),
	)
	fmt.Println(
		cmdEntry(
			"switch",
			"<name>",
			"Switch the active connection (alias: use)",
		),
	)
	fmt.Println(
		cmdEntry(
			"disconnect",
			"",
			"Clear the active connection (alias: clear, unset)",
		),
	)
	fmt.Println(
		cmdEntry(
			"remove",
			"--conn <name>",
			"Remove a saved connection (alias: delete)",
		),
	)
	fmt.Println(cmdEntry("ls", "", "List all configured connections"))
	fmt.Println(
		cmdEntry(
			"status",
			"",
			"Show the current active connection (alias: test)",
		),
	)
	fmt.Println()

	// ── QUERIES ───────────────────────────────────────────────────
	fmt.Println(sectionHeader("QUERIES"))
	fmt.Println(
		cmdEntry("add", "<name> [sql]", "Save a named query (alias: save)"),
	)
	fmt.Println(cmdEntry("run", "<name|id>", "Execute a saved query"))
	fmt.Println(
		cmdEntry("remove", "<name|id>", "Remove a saved query (alias: delete)"),
	)
	fmt.Println(
		cmdEntry(
			"list",
			"[queries]",
			"List saved queries for current connection",
		),
	)
	fmt.Println(cmdEntry("edit", "queries", "Edit saved queries in $EDITOR"))
	fmt.Println(cmdEntry("history", "", "Show query execution history"))
	fmt.Println()

	// ── DATABASE ──────────────────────────────────────────────────
	fmt.Println(sectionHeader("DATABASE"))
	fmt.Println(cmdEntry("query", "--table=<t>", "Run SQL against a table"))
	fmt.Println(
		cmdEntry(
			"tables",
			"[table]",
			"List tables or query one directly (alias: t, explore)",
		),
	)
	fmt.Println(
		"  remove      " + styles.Faint.Render(
			"Remove a saved query by name/id, or remove a connection entirely (alias: delete)",
		),
	)
	fmt.Println(
		"  run         " + styles.Faint.Render(
			"Run a saved query by name or id (alias: query)",
		),
	)
	fmt.Println(
		"  shell       " + styles.Faint.Render(
			"Interactive REPL for running queries (alias: repl)",
		),
	)
	fmt.Println(
		"  tables      " + styles.Faint.Render("List or query database tables"),
	)
	fmt.Println(
		"  explore     " + styles.Faint.Render("Explore database schema"),
	)
	fmt.Println(
		"  list        " + styles.Faint.Render("List connections or queries"),
	)
	fmt.Println(
		"  info        " + styles.Faint.Render(
			"Show tables or views in current connection",
		),
	)
	fmt.Println(
		"  edit        " + styles.Faint.Render(
			"Edit queries in your editor",
		),
	)
	fmt.Println(
		"  config      " + styles.Faint.Render(
			"Edit the main configuration file",
		),
	)
	fmt.Println(
		"  status      " + styles.Faint.Render(
			"Show the current active connection",
		),
	)
	fmt.Println(
		"  history     " + styles.Faint.Render(
			"Show query history (not implemented yet)",
		),
	)
	fmt.Println(
		"  explain     " + styles.Faint.Render(
			"Show relationships between tables",
		),
	)
	fmt.Println(
		"  help        " + styles.Faint.Render(
			"Show help for pam or a specific command",
		),
	)
	fmt.Println()

	// ── CONFIGURATION ─────────────────────────────────────────────
	fmt.Println(sectionHeader("CONFIGURATION"))
	fmt.Println(cmdEntry("edit", "config", "Edit the config file in $EDITOR"))
	fmt.Println(
		"  pam help              " + styles.Faint.Render("Show this help"),
	)
	fmt.Println(
		"  pam help <command>    " + styles.Faint.Render(
			"Show detailed help for a specific command",
		),
	)
	fmt.Println()

	// Examples
	fmt.Println(styles.Title.Render("Examples"))
	fmt.Println(
		"  pam init dev \"postgres://user:pass@localhost:5432/dbname\"",
	)
	fmt.Println(
		"  pam init oracle \"oracle://user:pass@localhost:1521/XEPDB1\"",
	)
	fmt.Println("  pam switch dev")
	fmt.Println("  pam add list_users \"SELECT * FROM users\"")
	fmt.Println("  pam run list_users")
	fmt.Println("  pam run \"select * from users\"")
	fmt.Println("  pam shell")
	fmt.Println("  pam list connections")
	fmt.Println("  pam list queries")
	fmt.Println("  pam edit config")
	fmt.Println("  pam edit queries")
}

func (a *App) PrintCommandHelp() {
	if len(os.Args) < 3 {
		a.PrintGeneralHelp()
		return
	}

	cmd := strings.ToLower(os.Args[2])

	section := func(title string) {
		fmt.Println(styles.Title.Render(title))
	}

	switch cmd {
	case "init", "create":
		section("Command:  init")
		fmt.Println(
			styles.Faint.Render(
				"Create and validate a new database connection. Launches interactive TUI if any required field is missing.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam init [flags]")
		fmt.Println(
			"  pam init <name> <connection-string>          # type auto-inferred",
		)
		fmt.Println(
			"  pam init <name> <db-type> <connection-string> [schema]  # legacy",
		)
		fmt.Println()
		section("Flags")
		fmt.Println("  --name,   -n          Connection name")
		fmt.Println(
			"  --type,   -t          Database type (optional, auto-inferred from conn string)",
		)
		fmt.Println(
			"  --conn,   -c          Connection string (alias: --conn-string)",
		)
		fmt.Println("  --schema, -s          Default schema (optional)")
		fmt.Println()
		section("Interactive TUI fields")
		fmt.Println(
			"  Connection name  — alias used to identify this connection",
		)
		fmt.Println(
			"  Database type    — postgres, mysql, sqlite, sqlserver, clickhouse, oracle, firebird",
		)
		fmt.Println(
			"  Host             — database host (or file path for SQLite)",
		)
		fmt.Println(
			"  Port             — auto-filled with the default port for the selected type",
		)
		fmt.Println("  Username         — database user")
		fmt.Println(
			"  Password         — masked by default, Ctrl+P to toggle visibility",
		)
		fmt.Println("  Database         — database / schema name")
		fmt.Println()
		fmt.Println(
			"  The connection string is assembled and previewed live as you type.",
		)
		fmt.Println()
		section("Examples")
		fmt.Println(
			"  pam init --name dev --conn \"postgres://user:pass@localhost:5432/dbname\"",
		)
		fmt.Println()
		fmt.Println("  # 2-arg positional with auto-inference")
		fmt.Println(
			"  pam init dev \"postgres://user:pass@localhost:5432/dbname\"",
		)
		fmt.Println()
		fmt.Println("  # 3-arg positional (legacy, explicit type)")
		fmt.Println(
			"  pam init dev postgres \"postgres://user:pass@localhost:5432/dbname\"",
		)
		fmt.Println()
		fmt.Println("  # Interactive mode")
		fmt.Println("  pam init")
		fmt.Println()
		fmt.Println("  # With schema")
		fmt.Println(
			"  pam init prod sqlserver \"sqlserver://sa:password@localhost:1433?database=mydb\" --schema public",
		)
		fmt.Println(
			"  pam init staging mysql \"user:pass@tcp(127.0.0.1:3306)/dbname\"",
		)
		fmt.Println()
		fmt.Println("  # DuckDB (included by default, requires CGO)")
		fmt.Println("  pam init local duckdb /path/to/mydb.db")

	case "switch", "use":
		section("Command: switch")
		fmt.Println(
			styles.Faint.Render(
				"Switch the active connection used by all other commands.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam switch/use <connection-name>")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  - Sets the connection to be used by 'add', 'run', 'list queries', etc.",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam switch dev")
		fmt.Println("  pam use prod")

	case "add", "save":
		section("Command: add")
		fmt.Println(
			styles.Faint.Render(
				"Save a new named query under the current connection.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam add <run-name> [query]")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  - If [query] is omitted, pam opens $EDITOR (default: vim) so you",
		)
		fmt.Println("    can write the query interactively.")
		fmt.Println("  - Each query gets a numeric ID as well as a name.")
		fmt.Println("  - Requires an active connection (use 'pam switch').")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam add list_users \"SELECT * FROM users\"")
		fmt.Println("  pam add update_status    # opens editor to write SQL")

	case "remove", "delete":
		section("Command: remove")
		fmt.Println(
			styles.Faint.Render(
				"Remove a saved query by name/id, or remove a connection entirely.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam remove <run-name-or-id>              # Remove query")
		fmt.Println(
			"  pam remove --connection <conn-name>    # Remove connection",
		)
		fmt.Println(
			"  pam remove -c <conn-name>             # Remove connection (short)",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam remove list_users                    # Remove query")
		fmt.Println(
			"  pam remove 3                             # Remove query by ID",
		)
		fmt.Println(
			"  pam remove --connection dev              # Remove connection",
		)
		fmt.Println(
			"  pam remove -c prod                         # Remove connection (short)",
		)

	case "query":
		section("Command: query")
		fmt.Println(
			styles.Faint.Render(
				"Run a SQL query against a table with primary key inference and interactive results.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println(
			"  pam run <query-name-or-id> [--edit | -e] [--last | -l] [--format | -f <fmt>]",
		)
		fmt.Println(
			"  pam run                      " + styles.Faint.Render(
				"# Opens the editor to build sql query",
			),
		)
		fmt.Println()
		section("Description")
		fmt.Println("  - Without SQL, defaults to 'SELECT * FROM <table>'.")
		fmt.Println(
			"  - Looks up a saved query by name or numeric ID and runs it against",
		)
		fmt.Println("    the current connection.")
		fmt.Println(
			"  - If no selector is provided, pam will open the editor to build sql query",
		)
		fmt.Println(
			"  - The result is rendered as an interactive table in your terminal.",
		)
		fmt.Println(
			"  - With '--edit' or '-e', pam opens the query in your $EDITOR before",
		)
		fmt.Println(
			"    running it and saves any changes back to the configuration.",
		)
		fmt.Println("  - With '--last' or '-l', runs the last used query")
		fmt.Println(
			"  - With '--format' or '-f', prints results to stdout instead of opening",
		)
		fmt.Println(
			"    the table UI. Formats: csv, json, tsv, html, sql, markdown",
		)
		fmt.Println()
		section("Interactive table view")
		fmt.Println(
			styles.Faint.Render(
				"When results are shown, you can interact with the table using the keyboard:",
			),
		)
		fmt.Println()
		fmt.Println(
			"  Arrow keys / h j k l  " + styles.Faint.Render(
				"Move selection around the table",
			),
		)
		fmt.Println(
			"  PageUp / Ctrl+u       " + styles.Faint.Render(
				"Scroll by a page up",
			),
		)
		fmt.Println(
			"  PageDown / Ctrl+d     " + styles.Faint.Render(
				"Scroll by a page down",
			),
		)
		fmt.Println(
			"  Home / 0 / _          " + styles.Faint.Render(
				"Jump to first row",
			),
		)
		fmt.Println(
			"  End / $               " + styles.Faint.Render(
				"Jump to last row",
			),
		)
		fmt.Println(
			"  g / G                 " + styles.Faint.Render(
				"Jump to top / bottom",
			),
		)
		fmt.Println(
			"  y / Enter             " + styles.Faint.Render(
				"Copy current cell value to clipboard (if supported)",
			),
		)
		fmt.Println(
			"  v                     " + styles.Faint.Render(
				"Start multi-selection mode",
			),
		)
		fmt.Println(
			"  u                     " + styles.Faint.Render(
				"Update selected cell",
			),
		)
		fmt.Println(
			"  d                     " + styles.Faint.Render(
				"Delete current row (requires WHERE clause)",
			),
		)
		fmt.Println(
			"  e                     " + styles.Faint.Render(
				"Open the editor to update and rerun query",
			),
		)
		fmt.Println(
			"  s                     " + styles.Faint.Render(
				"Save current query",
			),
		)
		fmt.Println(
			"  /                     " + styles.Faint.Render(
				"Search cell content",
			),
		)
		fmt.Println(
			"  n / N                 " + styles.Faint.Render(
				"Navigate to next/previous cell match",
			),
		)
		fmt.Println(
			"  f                     " + styles.Faint.Render(
				"Search column headers",
			),
		)
		fmt.Println(
			"  ; / ,                 " + styles.Faint.Render(
				"Navigate to next/previous column match",
			),
		)
		fmt.Println(
			"  Esc /Ctrl+c           " + styles.Faint.Render(
				"Quit the table view",
			),
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam run list_users")
		fmt.Println("  pam run \"select * from orders\"")
		fmt.Println("  pam run 2 --edit")
		fmt.Println("  pam run --last")
		fmt.Println("  pam run list_users -f json")
		fmt.Println(
			"  pam run \"SELECT * FROM users\" --format csv > users.csv",
		)
		fmt.Println("  pam query list_users")

	case "shell", "repl":
		section("Command: shell")
		fmt.Println(
			styles.Faint.Render(
				"Start an interactive REPL to run queries against the current connection.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam shell")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  - Opens REPL to run and list queries from the current active connection",
		)
		fmt.Println(
			"  - Supports inline SQL, saved queries by name/ID, and all run flags.",
		)
		fmt.Println(
			"  - Multi-line input: type SQL without trailing ; to continue.",
		)
		fmt.Println("  - Use up/down arrows to navigate command history.")
		fmt.Println()
		section("Meta-commands")
		fmt.Println("  exit, quit, \\q    Exit the REPL")
		fmt.Println("  help, \\h          Show help")
		fmt.Println("  list, ls, \\l      List saved queries or connections")
		fmt.Println("  status             Show connection info")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam shell")
		fmt.Println("  > select 1")
		fmt.Println("  > my-query")
		fmt.Println("  > my-query 123")
		fmt.Println("  > --last")
		fmt.Println("  > exit")

	case "list":
		section("Command: list")
		fmt.Println(styles.Faint.Render("List connections or saved queries."))
		fmt.Println()
		section("Usage")
		fmt.Println("  pam list [connections | queries] [search-term]")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  connections    List all configured connections (active one is highlighted)",
		)
		fmt.Println(
			"  queries        List saved queries for the current connection",
		)
		fmt.Println(
			"                 Optionally filter by search term (searches name and SQL)",
		)
		fmt.Println()
		section("Examples")
		fmt.Println(
			"  pam list                      # lists queries for the current connection",
		)
		fmt.Println("  pam list queries")
		fmt.Println(
			"  pam list queries emp          # list queries containing 'emp'",
		)
		fmt.Println(
			"  pam list queries employees    # list queries containing 'employees'",
		)
		fmt.Println(
			"  pam list queries --oneline    # list each query in one separate line",
		)
		fmt.Println("  pam list connections")

	case "tables":
		section("Command: tables")
		fmt.Println(
			styles.Faint.Render(
				"List all tables in the current database or query a specific table.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam tables [table-name] [--oneline | -o]")
		fmt.Println()
		section("Flags")
		fmt.Println(
			"  --oneline, -o    One table name per line (useful for scripting)",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam tables              # list all tables")
		fmt.Println("  pam tables users        # query the users table")
		fmt.Println("  pam tables --oneline    # list tables in oneline format")

	case "disconnect":
		section("Command: disconnect")
		fmt.Println(styles.Faint.Render("Clear the current active connection."))
		fmt.Println()
		section("Aliases")
		fmt.Println("  clear, unset")
		fmt.Println()
		section("Usage")
		fmt.Println("  pam disconnect")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  Clears the current active connection. You will need to use 'pam switch'",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam disconnect")

	case "edit":
		section("Command: edit")
		fmt.Println(
			styles.Faint.Render(
				"Edit queries in your editor.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam edit [<query-name-or-id>]")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  - Opens the editor to modify queries for the current connection.",
		)
		fmt.Println("    - With no arguments: opens all queries in one file")
		fmt.Println("    - With query name/id: edits a single query")
		fmt.Println(
			"    - Query name can be changed by editing the '-- queryname' header",
		)
		fmt.Println("  - Requires an active connection (use 'pam switch').")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam edit                    # edit all queries")
		fmt.Println("  pam edit list_users         # edit single query")
		fmt.Println("  pam edit 3                  # edit query by ID")
	case "config":
		section("Command: config")
		fmt.Println(
			styles.Faint.Render(
				"Edit the main configuration file.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam config")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  Opens the configuration file (~/.config/pam/config.yaml) in your editor.",
		)
		fmt.Println(
			"  Allows you to edit connections, color schemes, and other settings.",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam config")

	case "explore":
		section("Command: explore")
		fmt.Println(
			styles.Faint.Render(
				"Explore the database schema and query tables interactively.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam explore")
		fmt.Println("  pam explore <table> [--limit | -l N]")
		fmt.Println()
		section("Flags")
		fmt.Println(
			"  Without arguments, lists all tables and views in multi-column format.",
		)
		fmt.Println(
			"  With a table name, queries the table and shows results in an",
		)
		fmt.Println("  interactive table view (similar to 'pam run').")
		fmt.Println()
		fmt.Println(
			"  --limit, -l N  " + styles.Faint.Render(
				"Limit number of rows returned (default: from config or 1000)",
			),
		)
		fmt.Println()
		section("Examples")
		fmt.Println(
			"  pam explore                  # list all tables and views",
		)
		fmt.Println("  pam explore employees        # query employees table")
		fmt.Println("  pam explore orders -l 50     # query with 50 row limit")

	case "explain":
		section("Command: explain")
		fmt.Println(
			styles.Faint.Render(
				"Visualize foreign key relationships between tables.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam explain <table> [--depth | -d N]")
		fmt.Println()
		section("Flags")
		fmt.Println(
			"  --depth, -d N    Depth of relationships to traverse (default: 1)",
		)
		fmt.Println()
		section("Relationship types")
		fmt.Println(
			"  belongs to [N:1]    FK from this table pointing to another",
		)
		fmt.Println(
			"  has many   [1:N]    FK from another table pointing to this one",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam explain employees")
		fmt.Println("  pam explain employees --depth 2")
		fmt.Println("  pam explain departments -d 3")

	case "info":
		section("Command: info")
		fmt.Println(
			styles.Faint.Render(
				"Show all tables or views in the current connection.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam info <tables | views>")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam info tables")
		fmt.Println("  pam info views")

	case "status", "test":
		section("Command: status")
		fmt.Println(
			styles.Faint.Render(
				"Show the current active connection and test connectivity.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam status")

	case "history":
		section("Command: history")
		fmt.Println(styles.Faint.Render("Show query execution history."))
		fmt.Println()
		section("Usage")
		fmt.Println("  pam history")
		fmt.Println()
		section("Description")
		fmt.Println(
			styles.Faint.Render(
				"Export one or all tables from the active connection as a SQL dump.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam export [flags]")
		fmt.Println(
			"  pam export <table>           # shorthand for --table=<table>",
		)
		fmt.Println()
		section("Flags")
		fmt.Println("  --table,  -t <table>    Export only the specified table")
		fmt.Println(
			"  --output, -o <file>     Write dump to a file (default: stdout)",
		)
		fmt.Println("  --no-create             Skip CREATE TABLE statements")
		fmt.Println("  --drop                  Prepend DROP TABLE IF EXISTS")
		fmt.Println(
			"  --no-data               Schema only — skip INSERT statements (alias: --schema-only)",
		)
		fmt.Println(
			"  --data-only             Data only — skip CREATE TABLE (alias: --no-create)",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam export")
		fmt.Println("  pam export --table=users")
		fmt.Println("  pam export --output=dump.sql")
		fmt.Println("  pam export --drop --output=full.sql")
		fmt.Println("  pam export --no-data --output=schema.sql")
		fmt.Println("  pam export --data-only > inserts.sql")

	case "import":
		section("Command: import")
		fmt.Println(
			styles.Faint.Render(
				"Import a SQL dump file into the active connection.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam import [flags]")
		fmt.Println(
			"  pam import <file>                  # shorthand for --file=<file>",
		)
		fmt.Println("  cat dump.sql | pam import          # read from stdin")
		fmt.Println()
		section("Flags")
		fmt.Println("  --file,  -f <file>       SQL file to import")
		fmt.Println(
			"  --continue-on-error      Keep going after failed statements (alias: --continue)",
		)
		fmt.Println(
			"  --dry-run                Parse and list statements without executing",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam import dump.sql")
		fmt.Println("  pam import dump.sql --continue-on-error")
		fmt.Println("  pam import dump.sql --dry-run")
		fmt.Println("  cat dump.sql | pam import")

	case "completion":
		section("Command: completion")
		fmt.Println(
			styles.Faint.Render(
				"Generate and install shell completion scripts.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam completion <bash|zsh|fish> [--install]")
		fmt.Println()
		section("Flags")
		fmt.Println(
			"  --install    Write script to the standard shell path (auto-loads on new sessions)",
		)
		fmt.Println()
		section("Install paths")
		fmt.Println("  bash  → ~/.local/share/bash-completion/completions/pam")
		fmt.Println("  zsh   → ~/.zsh/completions/_pam")
		fmt.Println("  fish  → ~/.config/fish/completions/pam.fish")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam completion bash --install")
		fmt.Println("  pam completion zsh  --install")
		fmt.Println("  pam completion fish --install")
		fmt.Println("  pam completion bash              # print to stdout")
		fmt.Println(
			"  eval \"$(pam completion bash)\"    # load for current session only",
		)

	case "help":
		section("Command: help")
		fmt.Println(
			styles.Faint.Render(
				"Show general help or detailed help for a specific command.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam help [command]")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam help")
		fmt.Println("  pam help run")
		fmt.Println("  pam help list")

	default:
		fmt.Printf("%s  Unknown command %q — run %s for a list of commands.\n",
			styles.Error.Render("✗"),
			cmd,
			styles.Title.Render("pam help"),
		)
	}

	fmt.Println()
}

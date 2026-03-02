package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/caiolandgraf/pam/internal/styles"
	"github.com/charmbracelet/lipgloss"
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
	// Logo
	fmt.Println()
	fmt.Print(asciiLogo())
	fmt.Println()

	// Tagline + version
	tagline := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
		Render("Pam's database drawer — query manager for your databases")
	version := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Accent)).
		Bold(true).
		Render(Version)
	fmt.Printf("  %s   %s\n", tagline, version)
	fmt.Println()
	fmt.Println(separator())
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
		cmdEntry(
			"table-view",
			"<table>",
			"Inspect and edit table columns (alias: tv)",
		),
	)
	fmt.Println(
		cmdEntry(
			"info",
			"<tables|views>",
			"Show tables or views in current connection",
		),
	)
	fmt.Println(
		cmdEntry("explain", "<table>", "Visualize foreign key relationships"),
	)
	fmt.Println()

	// ── IMPORT / EXPORT ───────────────────────────────────────────
	fmt.Println(sectionHeader("IMPORT / EXPORT"))
	fmt.Println(
		cmdEntry("export", "[--table=<t>]", "Export tables as SQL dump"),
	)
	fmt.Println(
		cmdEntry(
			"import",
			"<file>",
			"Import a SQL dump into the active connection",
		),
	)
	fmt.Println()

	// ── CONFIGURATION ─────────────────────────────────────────────
	fmt.Println(sectionHeader("CONFIGURATION"))
	fmt.Println(cmdEntry("edit", "config", "Edit the config file in $EDITOR"))
	fmt.Println(
		cmdEntry(
			"completion",
			"<bash|zsh|fish>",
			"Generate shell completion script",
		),
	)
	fmt.Println()

	fmt.Println(separator())
	fmt.Println()

	// ── QUICK START ───────────────────────────────────────────────
	fmt.Println(sectionHeader("QUICK START"))
	accent := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Accent))
	faint := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Muted))
	fmt.Printf(
		"  %s  %s\n",
		accent.Render("1."),
		faint.Render("pam init                 # create your first connection"),
	)
	fmt.Printf(
		"  %s  %s\n",
		accent.Render("2."),
		faint.Render("pam t                    # list all tables"),
	)
	fmt.Printf(
		"  %s  %s\n",
		accent.Render("3."),
		faint.Render("pam t users              # query a table"),
	)
	fmt.Printf(
		"  %s  %s\n",
		accent.Render("4."),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Render("enjoy pam :D"),
	)
	fmt.Println()

	// ── HELP ──────────────────────────────────────────────────────
	fmt.Println(sectionHeader("HELP"))
	fmt.Printf(
		"  %s\n",
		faint.Render(
			"pam help <command>    show detailed help for any command",
		),
	)
	fmt.Println()
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
		fmt.Println(
			"  pam init                                                    # interactive TUI",
		)
		fmt.Println("  pam init [flags]")
		fmt.Println(
			"  pam init <name> <connection-string>                         # type auto-inferred",
		)
		fmt.Println(
			"  pam init <name> <db-type> <connection-string> [schema]      # explicit type",
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
			"  pam init                                                    # interactive",
		)
		fmt.Println(
			"  pam init --name dev --conn \"postgres://user:pass@localhost:5432/mydb\"",
		)
		fmt.Println(
			"  pam init dev \"postgres://user:pass@localhost:5432/mydb\"",
		)
		fmt.Println(
			"  pam init dev postgres \"postgres://user:pass@localhost:5432/mydb\"",
		)
		fmt.Println(
			"  pam init prod sqlserver \"sqlserver://sa:pass@localhost:1433?database=mydb\"",
		)

	case "switch", "use":
		section("Command: switch")
		fmt.Println(
			styles.Faint.Render(
				"Switch the active connection used by all other commands.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam switch <connection-name>")
		fmt.Println("  pam use    <connection-name>")
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
		fmt.Println("  pam add <name> [sql]")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  - If [sql] is omitted, opens $EDITOR (default: vim) to write the query.",
		)
		fmt.Println("  - Each query gets a numeric ID in addition to its name.")
		fmt.Println("  - Requires an active connection (pam switch <name>).")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam add list_users \"SELECT * FROM users\"")
		fmt.Println("  pam add update_status          # opens editor")

	case "remove", "delete":
		section("Command: remove")
		fmt.Println(
			styles.Faint.Render(
				"Remove a saved query by name or ID, or remove an entire connection.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println(
			"  pam remove <query-name-or-id>          # remove a saved query",
		)
		fmt.Println(
			"  pam remove --conn <connection-name>    # remove a connection",
		)
		fmt.Println(
			"  pam remove -c    <connection-name>     # same, short flag",
		)
		fmt.Println()
		section("Description")
		fmt.Println(
			"  - When removing a connection that is currently active, the active",
		)
		fmt.Println("    connection is cleared automatically.")
		fmt.Println(
			"  - Use 'pam switch <name>' to select another connection afterwards.",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam remove list_users")
		fmt.Println("  pam remove 3")
		fmt.Println("  pam remove --conn mydb")
		fmt.Println("  pam remove -c staging")
		fmt.Println("  pam delete --conn old_prod")

	case "query":
		section("Command: query")
		fmt.Println(
			styles.Faint.Render(
				"Run a SQL query against a table with primary key inference and interactive results.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam query --table=<table> [sql]")
		fmt.Println("  pam query --table=<table> --edit")
		fmt.Println()
		section("Flags")
		fmt.Println(
			"  --table, -t <name>    Target table (used for metadata and default SELECT *)",
		)
		fmt.Println(
			"  --edit,  -e           Open SQL in $EDITOR before executing",
		)
		fmt.Println()
		section("Description")
		fmt.Println("  - Without SQL, defaults to 'SELECT * FROM <table>'.")
		fmt.Println(
			"  - The --table flag enables row editing and deletion in the interactive view.",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam query --table=users")
		fmt.Println(
			"  pam query --table=users \"SELECT * FROM users WHERE active = 1\"",
		)
		fmt.Println("  pam query -t orders \"SELECT id, total FROM orders\"")
		fmt.Println("  pam query --table=users --edit")

	case "run":
		section("Command: run")
		fmt.Println(
			styles.Faint.Render(
				"Execute a saved query and display results in an interactive table view.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam run <name|id> [--edit | -e] [--last | -l]")
		fmt.Println(
			"  pam run                               # opens editor to build query on the fly",
		)
		fmt.Println()
		section("Flags")
		fmt.Println("  --edit, -e    Open the query in $EDITOR before running")
		fmt.Println("  --last, -l    Run the last used query")
		fmt.Println()
		section("Interactive table keys")
		fmt.Println("  ↑↓ / h j k l      Navigate rows and columns")
		fmt.Println("  PageUp / Ctrl+u   Scroll page up")
		fmt.Println("  PageDown / Ctrl+d Scroll page down")
		fmt.Println("  g / G             Jump to top / bottom")
		fmt.Println("  y / Enter         Copy current cell to clipboard")
		fmt.Println("  u                 Update selected cell")
		fmt.Println("  d                 Delete current row")
		fmt.Println("  e                 Edit and rerun query")
		fmt.Println("  s                 Save current query")
		fmt.Println("  Esc / Ctrl+c      Quit")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam run list_users")
		fmt.Println("  pam run 2 --edit")
		fmt.Println("  pam run --last")
		fmt.Println("  pam run \"SELECT * FROM orders\"")

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
		fmt.Println("  pam list")
		fmt.Println("  pam list queries")
		fmt.Println("  pam list queries users")
		fmt.Println("  pam list queries --oneline")
		fmt.Println("  pam list connections")
		fmt.Println(
			"  pam ls                           # shorthand for list connections",
		)

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
		fmt.Println("  pam tables")
		fmt.Println("  pam tables users")
		fmt.Println("  pam tables --oneline")

	case "table-view", "tv":
		section("Command: table-view")
		fmt.Println(
			styles.Faint.Render(
				"Inspect and edit the structure of a table — columns, types, constraints.",
			),
		)
		fmt.Println()
		section("Aliases")
		fmt.Println("  tv")
		fmt.Println()
		section("Usage")
		fmt.Println("  pam table-view <table>")
		fmt.Println("  pam tv         <table>")
		fmt.Println()
		section("Interactive keys")
		fmt.Println("  j / k  ↑↓     Navigate columns")
		fmt.Println("  a              Add a new column")
		fmt.Println("  e              Edit / alter selected column")
		fmt.Println("  r              Rename selected column")
		fmt.Println("  D              Drop selected column")
		fmt.Println("  q / Ctrl+c     Quit")
		fmt.Println()
		fmt.Println(
			"  Each action opens $EDITOR with a pre-filled SQL statement.",
		)
		fmt.Println("  Save and close to execute, clear content to cancel.")
		fmt.Println()
		section("Examples")
		fmt.Println("  pam table-view users")
		fmt.Println("  pam tv employees")

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
			"  Clears the active connection. Use 'pam switch <name>' to select one again.",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam disconnect")

	case "edit":
		section("Command: edit")
		fmt.Println(
			styles.Faint.Render(
				"Open pam's configuration or saved queries in $EDITOR.",
			),
		)
		fmt.Println()
		section("Usage")
		fmt.Println("  pam edit [config | queries]")
		fmt.Println()
		section("Description")
		fmt.Println(
			"  config    Edit the main config file (connections, settings, etc.)",
		)
		fmt.Println(
			"  queries   Edit all saved queries for the current connection",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam edit           # defaults to config")
		fmt.Println("  pam edit config")
		fmt.Println("  pam edit queries")

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
			"  --limit, -l N    Limit rows returned (default: config value or 1000)",
		)
		fmt.Println()
		section("Examples")
		fmt.Println("  pam explore")
		fmt.Println("  pam explore employees")
		fmt.Println("  pam explore orders -l 50")

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
		fmt.Println("  pam test")

	case "history":
		section("Command: history")
		fmt.Println(styles.Faint.Render("Show query execution history."))
		fmt.Println()
		section("Usage")
		fmt.Println("  pam history")

	case "export":
		section("Command: export")
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
		fmt.Println("  pam help remove")

	default:
		fmt.Printf("%s  Unknown command %q — run %s for a list of commands.\n",
			styles.Error.Render("✗"),
			cmd,
			styles.Title.Render("pam help"),
		)
	}

	fmt.Println()
}

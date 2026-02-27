package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eduardofuncao/pam/internal/config"
	"github.com/eduardofuncao/pam/internal/db"
	"github.com/eduardofuncao/pam/internal/spinner"
	"github.com/eduardofuncao/pam/internal/styles"
)

type exportFlags struct {
	tableName    string
	outputFile   string
	noCreate     bool
	dropIfExists bool
	noData       bool
}

func parseExportFlags() exportFlags {
	flags := exportFlags{}
	args := os.Args[2:]

	i := 0
	for i < len(args) {
		arg := args[i]

		switch {
		// --table / -t
		case arg == "--table" || arg == "-t":
			if i+1 < len(args) {
				flags.tableName = args[i+1]
				i += 2
			} else {
				printError("--table requires a value")
			}
		case strings.HasPrefix(arg, "--table="):
			flags.tableName = strings.TrimPrefix(arg, "--table=")
			i++
		case strings.HasPrefix(arg, "-t="):
			flags.tableName = strings.TrimPrefix(arg, "-t=")
			i++

		// --output / -o
		case arg == "--output" || arg == "-o":
			if i+1 < len(args) {
				flags.outputFile = args[i+1]
				i += 2
			} else {
				printError("--output requires a value")
			}
		case strings.HasPrefix(arg, "--output="):
			flags.outputFile = strings.TrimPrefix(arg, "--output=")
			i++
		case strings.HasPrefix(arg, "-o="):
			flags.outputFile = strings.TrimPrefix(arg, "-o=")
			i++

		// boolean flags
		case arg == "--no-create" || arg == "--no-create-table":
			flags.noCreate = true
			i++
		case arg == "--drop" || arg == "--drop-if-exists":
			flags.dropIfExists = true
			i++
		case arg == "--no-data" || arg == "--schema-only":
			flags.noData = true
			i++
		case arg == "--data-only":
			flags.noCreate = true
			i++

		default:
			// Treat a lone positional argument as the table name (convenience).
			if !strings.HasPrefix(arg, "-") && flags.tableName == "" {
				flags.tableName = arg
			}
			i++
		}
	}

	return flags
}

func (a *App) handleExport() {
	if a.config.CurrentConnection == "" {
		printError(
			"No active connection. Use 'pam switch <connection>' or 'pam init' first",
		)
	}

	flags := parseExportFlags()

	// Validate conflicting flags
	if flags.noCreate && flags.dropIfExists {
		printError("--no-create and --drop cannot be used together")
	}
	if flags.noCreate && flags.noData {
		printError("--no-create and --no-data would produce an empty dump")
	}

	// Determine output destination.
	// SQL always goes to out; progress/status messages always go to stderr
	// so they never pollute a redirect (pam export > dump.sql).
	var out *os.File
	if flags.outputFile != "" {
		f, err := os.Create(flags.outputFile)
		if err != nil {
			printError(
				"Could not create output file %q: %v",
				flags.outputFile,
				err,
			)
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	// Progress always on stderr
	progress := os.Stderr

	conn := config.FromConnectionYaml(
		a.config.Connections[a.config.CurrentConnection],
	)

	// Show spinner only when output goes to a file (stdout might be piped).
	var done chan struct{}
	if flags.outputFile != "" {
		done = make(chan struct{})
		go spinner.CircleWaitWithTimer(done)
	}

	if err := conn.Open(); err != nil {
		if done != nil {
			done <- struct{}{}
		}
		printError(
			"Could not open connection to %s/%s: %s",
			conn.GetDbType(),
			conn.GetName(),
			err,
		)
	}
	defer conn.Close()

	// Resolve the list of tables to export.
	var tables []string
	if flags.tableName != "" {
		tables = []string{flags.tableName}
	}
	// If tables is empty, ExportSQL will call GetTables() internally.

	opts := db.ExportOptions{
		IncludeCreate: !flags.noCreate,
		DropIfExists:  flags.dropIfExists,
		NoData:        flags.noData,
		Output:        out,
		Progress:      progress,
	}

	start := time.Now()

	if err := db.ExportSQL(conn, tables, opts); err != nil {
		if done != nil {
			done <- struct{}{}
		}
		printError("Export failed: %v", err)
	}

	elapsed := time.Since(start)

	if done != nil {
		done <- struct{}{}
	}

	// Print a summary to stderr (so it doesn't pollute a redirected file).
	if flags.outputFile != "" {
		fmt.Fprintf(
			os.Stderr,
			"%s Exported to %s in %s\n",
			styles.Success.Render("✓"),
			styles.Title.Render(flags.outputFile),
			elapsed.Round(time.Millisecond),
		)
	}
}

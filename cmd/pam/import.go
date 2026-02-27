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

type importFlags struct {
	inputFile       string
	continueOnError bool
	dryRun          bool
}

func parseImportFlags() importFlags {
	flags := importFlags{}
	args := os.Args[2:]

	i := 0
	for i < len(args) {
		arg := args[i]

		switch {
		// --file / -f
		case arg == "--file" || arg == "-f":
			if i+1 < len(args) {
				flags.inputFile = args[i+1]
				i += 2
			} else {
				printError("--file requires a value")
			}
		case strings.HasPrefix(arg, "--file="):
			flags.inputFile = strings.TrimPrefix(arg, "--file=")
			i++
		case strings.HasPrefix(arg, "-f="):
			flags.inputFile = strings.TrimPrefix(arg, "-f=")
			i++

		// --continue-on-error
		case arg == "--continue-on-error" || arg == "--continue":
			flags.continueOnError = true
			i++

		// --dry-run
		case arg == "--dry-run":
			flags.dryRun = true
			i++

		default:
			// Treat a lone positional argument as the input file (convenience).
			if !strings.HasPrefix(arg, "-") && flags.inputFile == "" {
				flags.inputFile = arg
			}
			i++
		}
	}

	return flags
}

func (a *App) handleImport() {
	if a.config.CurrentConnection == "" {
		printError(
			"No active connection. Use 'pam switch <connection>' or 'pam init' first",
		)
	}

	flags := parseImportFlags()

	// Determine input source: file or stdin.
	var input *os.File
	var inputName string

	if flags.inputFile != "" {
		f, err := os.Open(flags.inputFile)
		if err != nil {
			printError("Could not open input file %q: %v", flags.inputFile, err)
		}
		defer f.Close()
		input = f
		inputName = flags.inputFile
	} else {
		// Accept stdin only when it is a pipe/redirect, not an interactive TTY.
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			printError(
				"No input file specified.\n" +
					"  Use --file=<file> or pipe SQL via stdin:\n" +
					"    pam import --file=dump.sql\n" +
					"    cat dump.sql | pam import",
			)
		}
		input = os.Stdin
		inputName = "<stdin>"
	}

	conn := config.FromConnectionYaml(
		a.config.Connections[a.config.CurrentConnection],
	)

	// Spinner while connecting.
	done := make(chan struct{})
	go spinner.CircleWaitWithTimer(done)

	if err := conn.Open(); err != nil {
		done <- struct{}{}
		printError(
			"Could not open connection to %s/%s: %s",
			conn.GetDbType(),
			conn.GetName(),
			err,
		)
	}

	// Stop spinner before streaming progress lines.
	done <- struct{}{}

	defer conn.Close()

	fmt.Fprintf(
		os.Stderr,
		"Importing %s into %s/%s",
		styles.Title.Render(inputName),
		conn.GetDbType(),
		styles.Title.Render(conn.GetName()),
	)
	if flags.dryRun {
		fmt.Fprintf(os.Stderr, " %s", styles.Faint.Render("(dry-run)"))
	}
	fmt.Fprintln(os.Stderr)

	opts := db.ImportOptions{
		ContinueOnError: flags.continueOnError,
		DryRun:          flags.dryRun,
		Progress:        os.Stderr,
	}

	start := time.Now()
	result, err := db.ImportSQL(conn, input, opts)
	elapsed := time.Since(start)

	fmt.Fprintln(os.Stderr)

	// ── Summary ────────────────────────────────────────────────────────────

	if flags.dryRun {
		fmt.Fprintf(
			os.Stderr,
			"%s Dry run complete — %d statement(s) parsed from %s in %s\n",
			styles.Success.Render("✓"),
			result.Total,
			inputName,
			elapsed.Round(time.Millisecond),
		)
		return
	}

	// Hard stop (first-error mode).
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s Import aborted after %d/%d statement(s) in %s\n",
			styles.Error.Render("✗"),
			result.Executed,
			result.Total,
			elapsed.Round(time.Millisecond),
		)
		os.Exit(1)
	}

	// Completed (possibly with collected errors).
	if len(result.Errors) == 0 {
		fmt.Fprintf(
			os.Stderr,
			"%s Imported %d statement(s) from %s in %s\n",
			styles.Success.Render("✓"),
			result.Executed,
			inputName,
			elapsed.Round(time.Millisecond),
		)
	} else {
		fmt.Fprintf(
			os.Stderr,
			"%s Imported %d/%d statement(s) from %s in %s — %d error(s)\n",
			styles.Error.Render("!"),
			result.Executed,
			result.Total,
			inputName,
			elapsed.Round(time.Millisecond),
			len(result.Errors),
		)

		fmt.Fprintf(
			os.Stderr,
			"\n%s Failed statements:\n",
			styles.Error.Render("✗"),
		)
		for i, ie := range result.Errors {
			stmt := strings.ReplaceAll(ie.Statement, "\n", " ")
			if len([]rune(stmt)) > 120 {
				stmt = string([]rune(stmt)[:120]) + "..."
			}
			fmt.Fprintf(
				os.Stderr,
				"  [%d] %s\n      %s\n",
				i+1,
				styles.Error.Render(ie.Err.Error()),
				styles.Faint.Render(stmt),
			)
		}
	}
}

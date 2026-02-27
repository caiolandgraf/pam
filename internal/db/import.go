package db

import (
	"fmt"
	"io"
	"strings"
)

// ImportOptions configures how the SQL import is performed.
type ImportOptions struct {
	// ContinueOnError: if true, import continues after errors (collecting them).
	ContinueOnError bool
	// DryRun: if true, parse statements but do not execute them.
	DryRun bool
	// Progress is the destination writer for status messages.
	Progress io.Writer
}

// ImportError holds a single failed statement and its error.
type ImportError struct {
	Index     int
	Statement string
	Err       error
}

// ImportResult holds the summary of an import operation.
type ImportResult struct {
	Total    int
	Executed int
	Skipped  int // empty / whitespace-only statements
	Errors   []ImportError
}

// ImportSQL reads SQL from r, splits it into statements, and executes each one.
// It returns a summary result and a non-nil error only when the import is
// aborted early (ContinueOnError == false and a statement fails).
func ImportSQL(
	conn DatabaseConnection,
	r io.Reader,
	opts ImportOptions,
) (*ImportResult, error) {
	if opts.Progress == nil {
		opts.Progress = io.Discard
	}

	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	statements := SplitSQLStatements(string(raw))

	fmt.Fprintf(
		opts.Progress,
		"Parsed %d statement(s).\n",
		len(statements),
	)

	result := &ImportResult{}

	for idx, stmt := range statements {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			result.Skipped++
			continue
		}

		result.Total++

		if opts.DryRun {
			fmt.Fprintf(
				opts.Progress,
				"  [%d] (dry-run) %s\n",
				result.Total,
				truncateStmt(trimmed, 80),
			)
			result.Executed++
			continue
		}

		if err := conn.Exec(trimmed); err != nil {
			ie := ImportError{
				Index:     idx + 1,
				Statement: trimmed,
				Err:       err,
			}
			result.Errors = append(result.Errors, ie)

			fmt.Fprintf(
				opts.Progress,
				"  ✗ [%d] %v\n      SQL: %s\n",
				result.Total,
				err,
				truncateStmt(strings.ReplaceAll(trimmed, "\n", " "), 100),
			)

			if !opts.ContinueOnError {
				return result, fmt.Errorf(
					"import stopped at statement %d: %w",
					result.Total,
					err,
				)
			}
		} else {
			result.Executed++
		}
	}

	return result, nil
}

// SplitSQLStatements splits a SQL string into individual statements, correctly
// handling:
//   - single-quoted string literals (including ” escaped quotes)
//   - -- single-line comments
//   - /* */ block comments
//
// Each returned statement has its surrounding whitespace trimmed and does NOT
// include the trailing semicolon.
func SplitSQLStatements(sql string) []string {
	type lexState int
	const (
		stateNormal       lexState = iota
		stateInString              // inside '...'
		stateLineComment           // after --
		stateBlockComment          // inside /* ... */
	)

	var (
		statements []string
		current    strings.Builder
		state      = stateNormal
		runes      = []rune(sql)
		n          = len(runes)
	)

	for i := 0; i < n; i++ {
		ch := runes[i]

		switch state {

		case stateNormal:
			switch ch {
			case '\'':
				state = stateInString
				current.WriteRune(ch)

			case '-':
				if i+1 < n && runes[i+1] == '-' {
					state = stateLineComment
					i++ // consume second '-'
				} else {
					current.WriteRune(ch)
				}

			case '/':
				if i+1 < n && runes[i+1] == '*' {
					state = stateBlockComment
					i++ // consume '*'
				} else {
					current.WriteRune(ch)
				}

			case ';':
				// Statement boundary — emit if non-empty.
				if stmt := strings.TrimSpace(current.String()); stmt != "" {
					statements = append(statements, stmt)
				}
				current.Reset()

			default:
				current.WriteRune(ch)
			}

		case stateInString:
			current.WriteRune(ch)
			if ch == '\'' {
				// Two consecutive single quotes inside a string are an escape
				// sequence for a literal quote — keep both and stay in string.
				if i+1 < n && runes[i+1] == '\'' {
					current.WriteRune('\'')
					i++ // consume the second quote
				} else {
					// Closing quote — back to normal.
					state = stateNormal
				}
			}

		case stateLineComment:
			// Discard everything until end-of-line; the newline itself is
			// preserved so that line-number tracking remains meaningful.
			if ch == '\n' {
				state = stateNormal
				current.WriteRune(ch)
			}

		case stateBlockComment:
			// Discard everything until the closing '*/'.
			if ch == '*' && i+1 < n && runes[i+1] == '/' {
				state = stateNormal
				i++ // consume '/'
			}
		}
	}

	// Flush any trailing content that had no terminating semicolon.
	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		statements = append(statements, stmt)
	}

	return statements
}

// truncateStmt returns s truncated to maxLen runes, appending "..." when cut.
func truncateStmt(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

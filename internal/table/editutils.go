package table

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// buildEditorCommand creates an exec. Cmd with cursor positioning based on the editor type
func buildEditorCommand(editorCmd, tmpPath, content string, cursorHint cursorPositionHint) *exec.Cmd {
	line, col := findCursorPosition(content, cursorHint)

	switch editorCmd {
	case "vim", "nvim":
		// For vim/neovim:  +call cursor(line, col)
		return exec.Command(editorCmd, fmt.Sprintf("+call cursor(%d,%d)", line, col), tmpPath)
	case "nano":
		// For nano: +LINE,COLUMN
		return exec.Command(editorCmd, fmt.Sprintf("+%d,%d", line, col), tmpPath)
	case "emacs":
		// For emacs: +LINE: COLUMN
		return exec.Command(editorCmd, fmt.Sprintf("+%d:%d", line, col), tmpPath)
	case "code", "vscode":
		// For VS Code: --goto file:line:column --wait
		return exec.Command(editorCmd, "--goto", fmt.Sprintf("%s:%d:%d", tmpPath, line, col), "--wait")
	default:
		// Fallback:  just open the file
		return exec.Command(editorCmd, tmpPath)
	}
}

// cursorPositionHint tells the function where to look for the cursor position
type cursorPositionHint int

const (
	CursorAtUpdateValue cursorPositionHint = iota // Inside the value in UPDATE SET col = 'value'
	CursorAtWhereClause                            // Inside the value in WHERE col = 'value'
	CursorAtEndOfFile                              // At the end of the file
)

// findCursorPosition finds the appropriate cursor position based on the hint
func findCursorPosition(content string, hint cursorPositionHint) (line int, col int) {
	lines := strings.Split(content, "\n")

	switch hint {
	case CursorAtUpdateValue:
		// Look for:  SET column = 'value'
		re := regexp.MustCompile(`SET\s+\w+\s*=\s*'`)
		for i, lineText := range lines {
			match := re.FindStringIndex(lineText)
			if match != nil {
				// Position cursor right after the opening quote
				return i + 1, match[1] + 1
			}
		}
		return 3, 1 // fallback

	case CursorAtWhereClause:
		// Look for: WHERE column = 'value' and position inside the quotes
		re := regexp.MustCompile(`WHERE\s+\w+\s*=\s*'`)
		for i, lineText := range lines {
			match := re.FindStringIndex(lineText)
			if match != nil {
				// Position cursor right after the opening quote in WHERE clause
				return i + 1, match[1] + 1
			}
		}
		// Fallback: just find WHERE and position after it
		for i, lineText := range lines {
			if strings.Contains(strings.ToUpper(lineText), "WHERE") {
				idx := strings.Index(strings.ToUpper(lineText), "WHERE")
				if idx != -1 {
					return i + 1, idx + 7 // Position after "WHERE "
				}
			}
		}
		return len(lines), 1 // fallback to last line

	case CursorAtEndOfFile:
		// Position at the end of the last non-empty line
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.TrimSpace(lines[i]) != "" {
				return i + 1, len(lines[i]) + 1
			}
		}
		return len(lines), 1

	default:
		return 1, 1
	}
}

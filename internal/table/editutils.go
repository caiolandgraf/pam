package table

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// buildEditorCommand creates an exec.Cmd with cursor positioning based on the editor type
func buildEditorCommand(
	editorCmd, tmpPath, content string,
	cursorHint cursorPositionHint,
) *exec.Cmd {
	line, col := findCursorPosition(content, cursorHint)
	switch editorCmd {
	case "vim", "nvim":
		// For vim/neovim: +call cursor(line, col)
		return exec.Command(
			editorCmd,
			fmt.Sprintf("+call cursor(%d,%d)", line, col),
			tmpPath,
		)
	case "nano":
		// For nano: +LINE,COLUMN
		return exec.Command(
			editorCmd,
			fmt.Sprintf("+%d,%d", line, col),
			tmpPath,
		)
	case "emacs":
		// For emacs: +LINE:COLUMN
		return exec.Command(
			editorCmd,
			fmt.Sprintf("+%d:%d", line, col),
			tmpPath,
		)
	case "code", "vscode":
		// For VS Code: --goto file:line:column --wait
		return exec.Command(
			editorCmd,
			"--goto",
			fmt.Sprintf("%s:%d:%d", tmpPath, line, col),
			"--wait",
		)
	default:
		// Fallback: just open the file
		return exec.Command(editorCmd, tmpPath)
	}
}

type cursorPositionHint int

const (
	CursorAtUpdateValue cursorPositionHint = iota // Inside the value in UPDATE SET col = 'value'
	CursorAtWhereClause                           // Inside the value in WHERE col = 'value'
	CursorAtEndOfFile                             // At the end of the file
)

func findCursorPosition(
	content string,
	hint cursorPositionHint,
) (line int, col int) {
	lines := strings.Split(content, "\n")
	switch hint {
	case CursorAtUpdateValue:
		// Captura o valor apos SET col =, incluindo aspas
		re := regexp.MustCompile(`(?i)SET\s+\w+\s*=\s*([^;]+)`)
		for i, lineText := range lines {
			match := re.FindStringSubmatchIndex(lineText)
			if match != nil {
				raw := lineText[match[2]:match[3]]
				trimmed := strings.TrimLeft(raw, " \t")
				leading := len(raw) - len(trimmed)
				// 0-based index do primeiro char do valor (sem espacos)
				valStart := match[2] + leading
				if len(trimmed) > 0 &&
					(trimmed[0] == '\'' || trimmed[0] == '"') {
					// Posiciona apos a aspa de abertura (dentro do valor)
					return i + 1, valStart + 2
				}
				// Valor sem aspas: posiciona no primeiro char
				return i + 1, valStart + 1
			}
		}
		// Fallback: sem SET encontrado
		return 3, 1

	case CursorAtWhereClause:
		// Tenta encontrar WHERE com valor entre aspas
		reQuoted := regexp.MustCompile(`(?i)WHERE\s+\w+\s*=\s*(['"][^;]+)`)
		for i, lineText := range lines {
			match := reQuoted.FindStringSubmatchIndex(lineText)
			if match != nil {
				// match[2] e o indice da aspa de abertura (0-based)
				// Posiciona apos a aspa (dentro do valor)
				return i + 1, match[2] + 2
			}
		}
		// Fallback: WHERE sem aspas - posiciona logo apos "WHERE "
		for i, lineText := range lines {
			upper := strings.ToUpper(lineText)
			idx := strings.Index(upper, "WHERE")
			if idx != -1 {
				// +2: pula o espaco apos WHERE e converte para 1-based
				return i + 1, idx + len("WHERE") + 2
			}
		}
		// Ultimo fallback: sem WHERE
		return len(lines), 1

	case CursorAtEndOfFile:
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.TrimSpace(lines[i]) != "" {
				return i + 1, len(lines[i]) + 1
			}
		}
		return 1, 1

	default:
		return 1, 1
	}
}

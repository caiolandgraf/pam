package table

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// buildEditorCommand creates an exec. Cmd with cursor positioning based on the editor type
func buildEditorCommand(
	editorCmd, tmpPath, content string,
	cursorHint cursorPositionHint,
) *exec.Cmd {
	line, col := findCursorPosition(content, cursorHint)

	switch editorCmd {
	case "vim", "nvim":
		// For vim/neovim:  +call cursor(line, col)
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
		// For emacs: +LINE: COLUMN
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
		// Fallback:  just open the file
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
		// Regex para encontrar SET column = 'value' ou SET column = value (sem aspas)
		// (?i) para case-insensitive
		// Captura o nome da coluna e o valor (com ou sem aspas)
		re := regexp.MustCompile(
			`(?i)SET\s+(\w+)\s*=\s*('([^']*)'|"([^"]*)"|(\S+))`,
		)
		for i, lineText := range lines {
			match := re.FindStringSubmatchIndex(lineText)
			if match != nil {
				// match[0] é o início do "SET ...", match[1] é o fim
				// match[2] é o início do nome da coluna, match[3] é o fim
				// match[4] é o início do valor (com aspas), match[5] é o fim
				// match[6] é o início do valor (sem aspas), match[7] é o fim
				// match[8] é o início do valor (sem aspas), match[9] é o fim

				// Tenta encontrar o início do valor.
				// Prioriza o valor entre aspas simples, depois aspas duplas, depois sem aspas.
				valueStart := -1
				if match[6] != -1 { // Valor entre aspas simples
					valueStart = match[6] + 1 // Posição após a aspa simples de abertura
				} else if match[8] != -1 { // Valor entre aspas duplas
					valueStart = match[8] + 1 // Posição após a aspa dupla de abertura
				} else if match[10] != -1 { // Valor sem aspas
					valueStart = match[10] // Posição do início do valor
				}

				if valueStart != -1 {
					return i + 1, valueStart + 1 // +1 para converter de índice 0-based para 1-based
				}
			}
		}
		// Fallback mais robusto: se não encontrar SET, vai para o final da primeira linha
		if len(lines) > 0 {
			return 1, len(lines[0]) + 1
		}
		return 1, 1 // Último fallback

	case CursorAtWhereClause:
		// Regex para encontrar WHERE column = 'value' ou WHERE column = value (sem aspas)
		// (?i) para case-insensitive
		re := regexp.MustCompile(
			`(?i)WHERE\s+(\w+)\s*=\s*('([^']*)'|"([^"]*)"|(\S+))`,
		)
		for i, lineText := range lines {
			match := re.FindStringSubmatchIndex(lineText)
			if match != nil {
				valueStart := -1
				if match[6] != -1 { // Valor entre aspas simples
					valueStart = match[6] + 1
				} else if match[8] != -1 { // Valor entre aspas duplas
					valueStart = match[8] + 1
				} else if match[10] != -1 { // Valor sem aspas
					valueStart = match[10]
				}

				if valueStart != -1 {
					return i + 1, valueStart + 1
				}
			}
		}
		// Fallback: se não encontrar um valor específico, posiciona após a palavra WHERE
		for i, lineText := range lines {
			upperLine := strings.ToUpper(lineText)
			idx := strings.Index(upperLine, "WHERE")
			if idx != -1 {
				return i + 1, idx + len("WHERE") + 1 // Posição após "WHERE "
			}
		}
		return len(lines), 1 // Último fallback

	case CursorAtEndOfFile:
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.TrimSpace(lines[i]) != "" {
				return i + 1, len(lines[i]) + 1
			}
		}
		return 1, 1 // Se o arquivo estiver vazio ou só com espaços em branco

	default:
		return 1, 1
	}
}

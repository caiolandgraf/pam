package display

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	// Maximum width for the content (adjust based on your terminal)
	maxWidth = 80

	// Styles
	containerStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		MarginBottom(1).
		Width(maxWidth)

	queryNameStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		MarginBottom(1)

	keywordStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	commentStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("243")).
		Italic(true)

	// SQL keywords that should trigger line breaks
	breakKeywords = []string{
		"SELECT", "FROM", "WHERE", "JOIN", "LEFT JOIN", "RIGHT JOIN",
		"INNER JOIN", "OUTER JOIN", "ON", "AND", "OR", "ORDER BY",
		"GROUP BY", "HAVING", "LIMIT", "OFFSET", "UNION", "INSERT",
		"UPDATE", "DELETE", "SET", "VALUES",
	}

	// All SQL keywords to highlight
	allKeywords = append(breakKeywords, []string{
		"AS", "IN", "NOT", "NULL", "IS", "LIKE", "BETWEEN", "EXISTS",
		"CASE", "WHEN", "THEN", "ELSE", "END", "DISTINCT", "ALL",
	}...)
)

func formatSQL(query string) string {
	// Trim whitespace
	query = strings.TrimSpace(query)

	// Step 1: Protect comments and strings from processing
	protected := make(map[string]string)
	counter := 0

	// Protect line comments (-- comment)
	commentRegex := regexp.MustCompile(`--[^\n]*`)
	query = commentRegex.ReplaceAllStringFunc(query, func(match string) string {
		placeholder := fmt.Sprintf("__COMMENT_%d__", counter)
		protected[placeholder] = match // Store original, will style later
		counter++
		return placeholder
	})

	// Protect multi-line comments (/* comment */)
	multiCommentRegex := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	query = multiCommentRegex.ReplaceAllStringFunc(query, func(match string) string {
		placeholder := fmt.Sprintf("__MLCOMMENT_%d__", counter)
		protected[placeholder] = match
		counter++
		return placeholder
	})

	// Protect strings ('string' and "string")
	stringRegex := regexp.MustCompile(`'[^']*'|"[^"]*"`)
	query = stringRegex.ReplaceAllStringFunc(query, func(match string) string {
		placeholder := fmt.Sprintf("__STRING_%d__", counter)
		protected[placeholder] = match
		counter++
		return placeholder
	})

	// Step 2: Add line breaks before major keywords
	for _, keyword := range breakKeywords {
		// Use word boundaries to match whole keywords only
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
		query = pattern.ReplaceAllStringFunc(query, func(match string) string {
			upper := strings.ToUpper(match)
			return "\n" + upper
		})
	}

	// Step 3: Clean up whitespace and normalize
	lines := strings.Split(query, "\n")
	var normalized []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			normalized = append(normalized, line)
		}
	}
	query = strings.Join(normalized, "\n")

	// Step 4: Wrap long lines (accounting for border padding: 2 on each side + 2 for border = 6 total)
	wrapWidth := maxWidth - 6
	var wrappedLines []string
	for _, line := range strings.Split(query, "\n") {
		// Check if line starts with a keyword
		startsWithKeyword := false
		keywordPrefix := ""
		for _, kw := range breakKeywords {
			if strings.HasPrefix(strings.ToUpper(line), kw) {
				startsWithKeyword = true
				keywordPrefix = kw
				break
			}
		}

		if startsWithKeyword && len(line) > wrapWidth {
			// Split the line: keyword on first line, rest indented
			rest := strings.TrimSpace(line[len(keywordPrefix):])
			wrappedLines = append(wrappedLines, keywordPrefix)
			
			// Wrap the rest with indentation
			wrapped := wordwrap.String(rest, wrapWidth-2)
			for _, wl := range strings.Split(wrapped, "\n") {
				wrappedLines = append(wrappedLines, "  "+wl)
			}
		} else if len(line) > wrapWidth {
			// Regular line wrapping with indentation for continuation
			wrapped := wordwrap.String(line, wrapWidth)
			splitLines := strings.Split(wrapped, "\n")
			for i, wl := range splitLines {
				if i == 0 {
					wrappedLines = append(wrappedLines, wl)
				} else {
					wrappedLines = append(wrappedLines, "  "+wl)
				}
			}
		} else {
			wrappedLines = append(wrappedLines, line)
		}
	}
	query = strings.Join(wrappedLines, "\n")

	// Step 5: Highlight all keywords
	for _, keyword := range allKeywords {
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(keyword) + `\b`)
		query = pattern.ReplaceAllStringFunc(query, func(match string) string {
			return keywordStyle.Render(strings.ToUpper(match))
		})
	}

	// Step 6: Restore protected content and apply styles
	for placeholder, original := range protected {
		// Apply comment style to comments
		if strings.HasPrefix(placeholder, "__COMMENT_") || strings.HasPrefix(placeholder, "__MLCOMMENT_") {
			query = strings.ReplaceAll(query, placeholder, commentStyle.Render(original))
		} else {
			query = strings.ReplaceAll(query, placeholder, original)
		}
	}

	return query
}

func RenderQuery(name, query string) string {
	title := queryNameStyle.Render("◆ " + name)
	formattedSQL := formatSQL(query)
	content := lipgloss.JoinVertical(lipgloss.Left, title, formattedSQL)
	return containerStyle.Render(content)
}

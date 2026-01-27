package params

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/eduardofuncao/pam/internal/db"
)

// SubstituteParameters replaces :param|default or :param syntax with DB-specific placeholders
// Returns transformed SQL and ordered argument values for prepared statements
func SubstituteParameters(sql string, paramValues map[string]string, conn db.DatabaseConnection) (string, []any, error) {
	if len(paramValues) == 0 {
		return sql, []any{}, nil
	}

	// Find all :param|default or :param patterns in order
	// Regex matches :param|value or :param where value can be quoted with spaces
	re := regexp.MustCompile(`:(\w+)(?:\|('(?:[^'\\]|\\.)*'|(?:[^'\s\\]+)))?`)
	matches := re.FindAllStringSubmatchIndex(sql, -1)

	if len(matches) == 0 {
		return sql, []any{}, nil
	}

	// Build ordered list of parameter values based on occurrence in SQL
	var orderedValues []any
	paramIndex := make(map[string]int) // Maps param name to its index (1-based)
	currentIndex := 1

	// Process matches in order of first appearance
	for _, match := range matches {
		// match[2:4] = group 1 (param name)
		paramName := sql[match[2]:match[3]]

		// Only add each param once (in order of first appearance)
		if _, exists := paramIndex[paramName]; !exists {
			if value, ok := paramValues[paramName]; ok {
				paramIndex[paramName] = currentIndex
				orderedValues = append(orderedValues, value)
				currentIndex++
			} else {
				return "", nil, fmt.Errorf("missing value for parameter: %s", paramName)
			}
		}
	}

	// Now replace :param|default or :param with appropriate placeholders
	result := replaceParamPlaceholders(sql, conn, paramIndex)

	return result, orderedValues, nil
}

// replaceParamPlaceholders replaces all :param|default or :param with DB-specific placeholders
func replaceParamPlaceholders(sql string, conn db.DatabaseConnection, paramIndex map[string]int) string {
	// Use same regex as initial extraction to handle quoted strings
	re := regexp.MustCompile(`:(\w+)(?:\|('(?:[^'\\]|\\.)*'|(?:[^'\s\\]+)))?`)

	// Replace all matches at once using ReplaceAllStringFunc
	result := re.ReplaceAllStringFunc(sql, func(match string) string {
		// Extract param name
		paramName := strings.TrimPrefix(match, ":")
		if pipeIdx := strings.Index(paramName, "|"); pipeIdx != -1 {
			paramName = paramName[:pipeIdx]
		}

		// Get placeholder for this param
		if index, ok := paramIndex[paramName]; ok {
			return conn.GetPlaceholder(index)
		}

		return match
	})

	return result
}

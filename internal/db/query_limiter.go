package db

import (
	"fmt"
	"regexp"
	"strings"
)

// WrapWithLimit adds a LIMIT clause to a SELECT query if needed
func WrapWithLimit(sql string, limit int, dbType string) string {
	if limit <= 0 {
		return sql // No limit requested
	}
	
	trimmedSQL := strings.TrimSpace(sql)
	upperSQL := strings.ToUpper(trimmedSQL)
	
	// Don't wrap non-SELECT queries
	if ! isSelectQuery(upperSQL) {
		return sql
	}
	
	// Check if query already has a limit
	if hasLimitClause(upperSQL, dbType) {
		return sql
	}
	
	return addLimitClause(trimmedSQL, limit, dbType)
}

func isSelectQuery(upperSQL string) bool {
	selectKeywords := []string{"SELECT", "WITH", "SHOW", "DESCRIBE", "DESC", "EXPLAIN", "PRAGMA"}
	for _, keyword := range selectKeywords {
		if strings.HasPrefix(upperSQL, keyword+" ") || upperSQL == keyword {
			return true
		}
	}
	return false
}

func hasLimitClause(upperSQL string, dbType string) bool {
	switch dbType {
	case "postgres", "postgresql", "mysql", "mariadb", "sqlite", "sqlite3":
		return regexp.MustCompile(`\bLIMIT\s+\d+`).MatchString(upperSQL)
		
	case "oracle", "godror":
		hasFetch := regexp.MustCompile(`\bFETCH\s+FIRST\s+\d+`).MatchString(upperSQL)
		hasRowNum := regexp.MustCompile(`\bROWNUM\s*[<>=]`).MatchString(upperSQL)
		return hasFetch || hasRowNum
		
	default: 
		return false
	}
}

func addLimitClause(sql string, limit int, dbType string) string {
	switch dbType {
	case "postgres", "postgresql": 
		return fmt.Sprintf("%s\nLIMIT %d", sql, limit)
	case "mysql", "mariadb":
		return fmt.Sprintf("%s\nLIMIT %d", sql, limit)
	case "sqlite", "sqlite3":
		return fmt.Sprintf("%s\nLIMIT %d", sql, limit)
	case "oracle", "godror": 
		return fmt.Sprintf("%s\nFETCH FIRST %d ROWS ONLY", sql, limit)
	default:
		return sql
	}
}

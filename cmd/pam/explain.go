package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eduardofuncao/pam/internal/config"
	"github.com/eduardofuncao/pam/internal/db"
	"github.com/eduardofuncao/pam/internal/styles"
)

type explainFlags struct {
	depth       int
	showColumns bool
}

func parseExplainFlags() (explainFlags, []string) {
	flags := explainFlags{
		depth:       3,
		showColumns: false,
	}
	remainingArgs := []string{}
	args := os.Args[2:]

	for i, arg := range args {
		if arg == "--depth" || arg == "-d" {
			if i+1 < len(args) {
				if depth, err := strconv.Atoi(args[i+1]); err == nil {
					flags.depth = depth
				}
			}
		} else if arg == "--columns" || arg == "-c" {
			flags.showColumns = true
		} else if !strings.HasPrefix(arg, "-") {
			remainingArgs = append(remainingArgs, arg)
		}
	}

	return flags, remainingArgs
}

func (a *App) handleExplain() {
	if a.config.CurrentConnection == "" {
		printError(
			"No active connection. Use 'pam switch <connection>' or 'pam init' first",
		)
	}

	flags, args := parseExplainFlags()

	if len(args) == 0 {
		fmt.Println("Usage: pam explain [--depth|-d N] [--columns|-c] <table-name>")
		os.Exit(1)
	}

	conn := config.FromConnectionYaml(
		a.config.Connections[a.config.CurrentConnection],
	)

	if err := conn.Open(); err != nil {
		printError(
			"Could not open connection to %s: %v",
			a.config.CurrentConnection,
			err,
		)
	}
	defer conn.Close()

	tableName := args[0]
	visited := make(map[string]bool)
	tree := a.buildRelationshipTree(conn, tableName, flags.showColumns, flags.depth, 0, visited)

	fmt.Println(tree)
}

type relationshipNode struct {
	tableName      string
	relationships  []relationship
}

type relationship struct {
	column           string
	referencedTable  string
	referencedColumn string
}

func (a *App) buildRelationshipTree(conn db.DatabaseConnection, tableName string, showColumns bool, maxDepth, currentDepth int, visited map[string]bool) string {
	if currentDepth >= maxDepth {
		return ""
	}

	node := &relationshipNode{
		tableName: tableName,
	}

	fks, err := conn.GetForeignKeys(tableName)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}

	for _, fk := range fks {
		node.relationships = append(node.relationships, relationship{
			column:           fk.Column,
			referencedTable:  fk.ReferencedTable,
			referencedColumn: fk.ReferencedColumn,
		})
	}

	return a.renderNode(conn, node, showColumns, maxDepth, currentDepth, visited)
}

func (a *App) renderNode(conn db.DatabaseConnection, node *relationshipNode, showColumns bool, maxDepth, currentDepth int, visited map[string]bool) string {
	var builder strings.Builder

	if currentDepth == 0 {
		metadata, _ := conn.GetTableMetadata(node.tableName)
		if len(metadata.PrimaryKeys) > 0 {
			pks := strings.Join(metadata.PrimaryKeys, ", ")
			builder.WriteString(styles.TableName.Render(node.tableName))
			builder.WriteString(" ")
			builder.WriteString(styles.PrimaryKeyLabel.Render(fmt.Sprintf("[%s]", pks)))
		} else {
			builder.WriteString(styles.TableName.Render(node.tableName))
		}
		builder.WriteString("\n")
	} else {
		builder.WriteString(styles.TableName.Render(node.tableName) + "\n")
	}

	visited[node.tableName] = true

	for i, rel := range node.relationships {
		isLast := i == len(node.relationships)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}

		var columnInfo string
		if showColumns {
			columnInfo = fmt.Sprintf("%s → ", rel.column)
		}

		cardinality := "[N:1]"
		relStyle := styles.BelongsToStyle

		builder.WriteString(styles.TreeConnector.Render(prefix))
		builder.WriteString(columnInfo)
		builder.WriteString(relStyle.Render("belongs to →"))
		builder.WriteString(" ")
		builder.WriteString(styles.CardinalityStyle.Render(cardinality))
		builder.WriteString(" ")
		builder.WriteString(styles.TableName.Render(rel.referencedTable))
		builder.WriteString("\n")

		isSelfReference := (rel.referencedTable == node.tableName)

		if isSelfReference {
			continue
		}

		if visited[rel.referencedTable] {
			continue
		}

		childNode := &relationshipNode{
			tableName: rel.referencedTable,
		}

		childFks, err := conn.GetForeignKeys(rel.referencedTable)
		if err != nil {
			continue
		}

		for _, fk := range childFks {
			childNode.relationships = append(childNode.relationships, relationship{
				column:           fk.Column,
				referencedTable:  fk.ReferencedTable,
				referencedColumn: fk.ReferencedColumn,
			})
		}

		childPrefix := "    "
		if !isLast {
			childPrefix = "│   "
		}

		childTree := a.renderNode(conn, childNode, showColumns, maxDepth, currentDepth+1, visited)
		lines := strings.Split(childTree, "\n")

		for j, line := range lines {
			if line == "" {
				continue
			}
			if j > 0 {
				builder.WriteString(styles.TreeConnector.Render(childPrefix))
			}
			builder.WriteString(line + "\n")
		}
	}

	return builder.String()
}

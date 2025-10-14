package db

import (
	"database/sql"
	"fmt"
	"log"
)

type Connection struct {
	Name       string
	DBType     string
	ConnString string
	Username   string
	Password   string
	DB         *sql.DB
	Queries    map[string]Query
}

type Query struct {
	Name string
	SQL  string
}

func NewConnection(name, dbType, connStr, user, pass string) *Connection {
	return &Connection{
		Name:       name,
		DBType:     dbType,
		ConnString: connStr,
		Username:   user,
		Password:   pass,
	}
}

func (c *Connection) Open() error {
	db, err := sql.Open(c.DBType, c.ConnString)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("ping db: %w", err)
	}
	c.DB = db
	return nil
}

func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

func (c *Connection) Query(queryName string, args ...interface{}) ([]string, [][]string, error) {
	rows, err := c.DB.Query(c.Queries[queryName].SQL, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Error getting columns: %v", err)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	var data [][]string
	for rows.Next() {
		err = rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		rowData := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				rowData[i] = "NULL"
			} else {
				rowData[i] = fmt.Sprintf("%v", val)
			}
		}
		data = append(data, rowData)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error during iteration: %v", err)
	}
	return columns, data, nil
}

package engine

import (
	"fmt"
	"strings"
	"sync"
)


type Database struct {
	tables map[string]*Table
	mu sync.RWMutex
}


func NewDatabase() *Database {
	return  &Database{
		tables: make(map[string]*Table),
	}
}



func (db *Database) CreateTable(name string, columns []Column)  (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.tables[name]; exists {
		return "", fmt.Errorf("table %s already exists", name)
	}

	table := NewTable(name, columns)
	db.tables[name] = table

	return  fmt.Sprintf("Table '%s' creates successfully", name), nil
}


func (db *Database) Insert(tableName string, values []interface{}) (string, error) {
		db.mu.Lock()
		table, exists := db.tables[tableName]
		db.mu.RUnlock()
		
		if !exists {
			return "", fmt.Errorf("table %s does not exist", tableName)
		}

		id, err := table.Insert(values)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Inserted row with ID: %d\n", id), nil
}

func (db *Database) Select(tableName string, columns []string, condition *Condition) (string, error) {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("table %s does not exist", tableName)
	}

	rows, err := table.Select(columns, condition)
	if err != nil {
		return "", err
	}

	if len(rows) == 0 {
		return "No rows found\n", nil
	}

	var result strings.Builder
	// Header
	if len(columns) == 1 && columns[0] == "*" {
		columns = nil
		if len(table.Columns) > 0 {
			for _, col := range table.Columns {
				columns = append(columns, col.Name)
			}
		}
	}

	result.WriteString("| ")
	for _, col := range columns {
		result.WriteString(fmt.Sprintf("%-15s | ", col))
	}
	result.WriteString("\n")
	result.WriteString(strings.Repeat("-", 20*len(columns)+3))
	result.WriteString("\n")

	// Rows
	for _, row := range rows {
		result.WriteString("| ")
		for _, col := range columns {
			val := row.Get(col)
			result.WriteString(fmt.Sprintf("%-15v | ", val))
		}
		result.WriteString("\n")
	}
	result.WriteString(fmt.Sprintf("\n%d row(s) returned\n", len(rows)))
	return result.String(), nil
}


func (db *Database) Update(tableName string, setValues map[string]interface{}, condition *Condition) (string, error) {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	db.mu.RUnlock()

	if !exists {
		return  "", fmt.Errorf("table %s does not exist", tableName)
	}
	updated, err := table.Update(setValues, condition)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Updated %d row(s)\n", updated), nil

}

func (db *Database) Delete(tableName string, condition *Condition) (string, error) {
	db.mu.RLock()
	table, exists := db.tables[tableName]
	if !exists {
		return "", fmt.Errorf("table %s does not exist", tableName)
	}

	deleted, err := table.Delete(condition)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Deleted %d row(s)\n", deleted), nil
}

func (db *Database) Close() error {
	return nil
}


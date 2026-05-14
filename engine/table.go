package engine

import (
	"fmt"
	"sync"
)

type Table struct {
	Name    string
	Columns []Column
	Rows    map[int64]*Row
	NextID  int64
	mu      sync.RWMutex
}

type Column struct {
	Name string
	Type string
}

func NewTable(name string, columns []Column) *Table {
	return &Table{
		Name:    name,
		Columns: columns,
		Rows:    make(map[int64]*Row),
		NextID:  1,
	}
}

func (t *Table) Insert(values []interface{}) (int64, error) {
	if len(values) != len(t.Columns) {
		return 0, fmt.Errorf("expected %d values, got %d", len(t.Columns), len(values))
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	row := NewRow()
	row.ID = t.NextID

	for i, col := range t.Columns {
		if err := t.validateType(col.Type, values[i]); err != nil {
			return 0, err
		}
		row.Set(col.Name, values[i])
	}

	t.Rows[t.NextID] = row
	t.NextID++

	return row.ID, nil
}

func (t *Table) Select(columns []string, condition *Condition) ([]*Row, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []*Row
	for _, row := range t.Rows {
		if condition != nil {
			matches, err := t.evaluateCondition(row, condition)
			if err != nil {
				return nil, err
			}
			if !matches {
				continue
			}
		}
		if len(columns) == 1 && columns[0] == "*" {
			result = append(result, row)
		} else {
			filteredRow := NewRow()
			filteredRow.ID = row.ID
			for _, col := range columns {
				if val, ok := row.Values[col]; ok {
					filteredRow.Set(col, val)
				}
			}
			result = append(result, filteredRow)
		}
	}
	return result, nil
}

func (t *Table) Update(setValues map[string]interface{}, condition *Condition) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	updated := 0
	for id, row := range t.Rows {
		if condition != nil {
			matches, err := t.evaluateCondition(row, condition)
			if err != nil {
				return 0, err
			}
			if !matches {
				continue
			}
		}
		// Update values
		for col, val := range setValues {
			var colType string
			for _, c := range t.columns {
				if c.Name == col {
					colType = c.Type
					break
				}
			}
			if err := t.validateType(colType, val); err != nil {
				return 0, err
			}
			row.Set(col, val)
		}
		updated++
	}
	return updated, nil
}

func (t *Table) Delete(condition *Condition) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	deleted := 0
	idsToDeelte := []int64{}

	for id, row := range t.Rows {
		if condition != nil {
			matches, err := t.evaluateCondition(row, condition)
			if err != nil {
				return 0, err
			}
			if !matches {
				continue
			}
		}
		idsToDeelte = append(idsToDeelte, id)
		deleted++
	}
	for _, id := range idsToDeelte {
		delete(t.Rows, id)
	}
	return deleted, nil
}


func (t *Table) validateType(colType string, val interface{}) error {
	switch colType {
	case "int":
		switch val.(type){
		case int, int64, float64:
			return nil
		default:
			return fmt.Errorf("expected int for column, got %T", val)
		}
	case "string":
		if _, ok := val.(string); !ok {
			return fmt.Errorf("expected string for column, got %T", val)
		}
		return nil
	default:
		return  fmt.Errorf("unknown column type: %s", colType)
	}
}


func (t *Table) evaluateCondition(row *Row, cond *Condition) (bool, error) {
    val, ok := row.Values[cond.Column]
    if !ok {
        return false, fmt.Errorf("column %s not found", cond.Column)
    }
    
    switch cond.Operator {
    case "=":
        return fmt.Sprintf("%v", val) == fmt.Sprintf("%v", cond.Value), nil
    case "!=":
        return fmt.Sprintf("%v", val) != fmt.Sprintf("%v", cond.Value), nil
    case ">":
        return compareNumbers(val, cond.Value) > 0, nil
    case "<":
        return compareNumbers(val, cond.Value) < 0, nil
    case ">=":
        return compareNumbers(val, cond.Value) >= 0, nil
    case "<=":
        return compareNumbers(val, cond.Value) <= 0, nil
    default:
        return false, fmt.Errorf("unknown operator: %s", cond.Operator)
    }
}

func compareNumbers(a, b interface{}) int {
    aFloat := toFloat64(a)
    bFloat := toFloat64(b)
    if aFloat < bFloat {
        return -1
    } else if aFloat > bFloat {
        return 1
    }
    return 0
}

func toFloat64(val interface{}) float64 {
    switch v := val.(type) {
    case int:
        return float64(v)
    case int64:
        return float64(v)
    case float64:
        return v
    default:
        return 0
    }
}
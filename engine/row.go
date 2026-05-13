package engine

import (
	"fmt"
	"reflect"
)

type Row struct {
	ID     int64
	Values map[string]interface{}
}

func NewRow() *Row {
	return &Row{
		Values: make(map[string]interface{}),
	}
}

func (r *Row) Set(col string, val interface{}) {
	r.Values[col] = val
}

func (r *Row) Get(col string) interface{} {
	return r.Values[col]
}

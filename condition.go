package mquery

import (
	"fmt"
	"html"
	"strings"
)

type conditionQuery struct {
	table tableQuery
	and   []string
	or    []string
}

type ConditionOption func(c *conditionQuery)

func Like(column interface{}, value interface{}) ConditionOption {
	return func(wb *conditionQuery) {
		colStr := fmt.Sprintf("%v", column)
		wb.table.colValid(colStr)
		valuee := fmt.Sprintf("%v", value)
		if valuee == "" {
			return
		}
		valuee = "'%" + html.EscapeString(valuee) + "%'"
		wb.and = append(wb.and, fmt.Sprintf("%s LIKE %s", colStr, valuee))
	}
}

func In(selectIn string, col ...interface{}) ConditionOption {
	return func(wb *conditionQuery) {
		colStr := make([]string, 0, len(col))
		for _, v := range col {
			column := fmt.Sprintf("%v", v)
			wb.table.colValid(column)
			colStr = append(colStr, column)
		}
		wb.and = append(wb.and, fmt.Sprintf("(%s) IN (%s)", strings.Join(colStr, ","), selectIn))
	}
}

type Operator string

const (
	EqualOps            Operator = "="
	LessOps             Operator = "<"
	GreaterOps          Operator = ">"
	LessThanEqualOps    Operator = "<="
	GreaterThanEqualOps Operator = ">="
	NotEqualOps         Operator = "<>"
)

func And(field interface{}, ops Operator, value ...interface{}) ConditionOption {
	return func(wb *conditionQuery) {
		wb.table.colValid(field)
		var v interface{}
		if len(value) > 0 {
			v = value[0]
		}
		wb.and = append(wb.and, fmt.Sprintf("%s %s %s", field, ops, interfaceToString(v)))
	}
}
func Or(field interface{}, ops Operator, value ...interface{}) ConditionOption {
	return func(wb *conditionQuery) {
		wb.table.colValid(field)
		var v interface{}
		if len(value) > 0 {
			v = value[0]
		}
		wb.or = append(wb.or, fmt.Sprintf("%s %s %s", field, ops, interfaceToString(v)))
	}
}

func newCondition(table tableQuery, options ...ConditionOption) toQuery {
	con := &conditionQuery{table: table}
	for _, setter := range options {
		setter(con)
	}
	return con
}
func (c conditionQuery) ToQuery() string {
	var con []string
	if c.and != nil {
		con = append(con, strings.Join(c.and, " AND "))
	}
	if c.or != nil {
		con = append(con, strings.Join(c.or, " OR "))
	}
	if con == nil {
		return ""
	}
	return strings.Join(con, " OR ")
}

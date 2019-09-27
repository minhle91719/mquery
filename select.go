package mquery

import (
	"fmt"
	"strings"
)

type selectQueryBuild struct {
	table   tableQuery
	field   []string
	asMap   map[string]string
	isCount bool
}

func (s selectQueryBuild) Where(opts ...WhereOption) toQuery {
	return newWhereQuery(s.table, s.ToQuery(), opts)
}

type SelectOption func(sq *selectQueryBuild)

func SelectField(list ...interface{}) SelectOption {
	return func(sq *selectQueryBuild) {
		for _, v := range list {
			colStr := fmt.Sprintf("%v", v)
			if colStr == "*" {
				sq.field = sq.table.column
				return
			}
			sq.table.colValid(colStr)
			sq.field = append(sq.field, colStr)
		}
	}
}
func SelectAll() SelectOption {
	return func(sq *selectQueryBuild) {
		sq.field = sq.table.column
	}
}

func SelectAs(column, as string) SelectOption {
	return func(sq *selectQueryBuild) {
		sq.table.colValid(column)
		sq.asMap[column] = as
	}
}
func Count() SelectOption {
	return func(sq *selectQueryBuild) {
		sq.isCount = true
	}
}

func (s *selectQueryBuild) ToQuery() string {
	var value = make([]string, 0, len(s.asMap)+len(s.field))
	if s.isCount {
		value = append(value, "COUNT(1)")
	} else {
		for _, v := range s.field {
			value = append(value, v)
		}
		for k, v := range s.asMap {
			value = append(value, fmt.Sprintf("%s AS %s", k, v))
		}
	}
	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(value, ","), s.table.tableName)
	return query
}

func newSelect(tb tableQuery, options []SelectOption) *selectQueryBuild {
	q := &selectQueryBuild{table: tb, asMap: make(map[string]string)}
	for _, setter := range options {
		setter(q)
	}
	return q
}

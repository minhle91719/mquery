package mquery

import (
	"fmt"
	"strings"
)

type insert struct {
	table    tableQuery
	field    []string
	isIgnore bool
	isValues bool
	rows     int64
	
	notCheckField bool
}

func (iq insert) ToQuery() string {
	var (
		valueQuery = "VALUE"
		insertType = "INSERT"
	)
	if iq.isIgnore {
		insertType = "INSERT IGNORE"
	}
	value := genValueParam(len(iq.field))
	var values = make([]string, 0, iq.rows)
	values = append(values, value)
	if iq.isValues {
		valueQuery = "VALUES"
		for i := int64(0); i < iq.rows-1; i++ {
			values = append(values, value)
		}
	}
	q := fmt.Sprintf("%s INTO %s(%s) %s%s", insertType, iq.table.tableName, strings.Join(iq.field, ","), valueQuery, strings.Join(values, ","))
	if iq.table.isLogger {
		iq.table.logger.Infof(q)
	}
	return q
}

type InsertOption func(i *insert)

func NotCheckFieldInsert() InsertOption {
	return func(i *insert) {
		i.notCheckField = true
	}
}
func WithField(field ...interface{}) InsertOption {
	return func(i *insert) {
		for _, v := range field {
			f := fmt.Sprintf("%v", v)
			if !i.notCheckField {
				i.table.colValid(f)
			}
			i.field = append(i.field, f)
		}
	}
}
func WithValues(rows int64) InsertOption {
	return func(i *insert) {
		i.isValues = true
		i.rows = rows
	}
}

func WithIgnore() InsertOption {
	return func(i *insert) {
		i.isIgnore = true
	}
}

func newInsert(table tableQuery, options []InsertOption) toQuery {
	i := &insert{table: table}
	for _, setter := range options {
		setter(i)
	}
	return i
}

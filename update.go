package mquery

import (
	"fmt"
	"strings"
)

type updateQueryBuild struct {
	table         tableQuery
	col           []string
	field         map[string]string
	notCheckField bool
}

type UpdateOption func(uq *updateQueryBuild)

func NotCheckFieldUpdate() UpdateOption {
	return func(uq *updateQueryBuild) {
		uq.notCheckField = true
	}
}
func UpdateField(colName interface{}, value interface{}) UpdateOption {
	return func(uq *updateQueryBuild) {
		if !uq.notCheckField {
			uq.table.colValid(colName)
		}
		column := fmt.Sprintf("%v", colName)
		if value == nil {
			uq.col = append(uq.col, column)
		} else {
			uq.field[column] = interfaceToString(value)
		}
	}
}
func (u *updateQueryBuild) ToQuery() string {
	var values = make([]string, 0, len(u.field))
	for _, v := range u.col {
		values = append(values, fmt.Sprintf("%s = ?", v))
	}
	for k, v := range u.field {
		values = append(values, fmt.Sprintf("%s = %s", k, v))
	}
	return fmt.Sprintf("UPDATE %s SET %s", u.table.tableName, strings.Join(values, ","))
}
func newUpdateQuery(tb tableQuery, opts []UpdateOption) *updateQueryBuild {
	u := &updateQueryBuild{table: tb, field: make(map[string]string)}
	for _, setter := range opts {
		setter(u)
	}
	return u
}
func (u *updateQueryBuild) Where(opts ...WhereOption) toQuery {
	return newWhereQuery(u.table, u.ToQuery(), opts)
}

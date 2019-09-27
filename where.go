package mquery

import (
	"fmt"
	"strings"
)

type whereQueryBuild struct {
	table tableQuery
	
	prefix    string
	condition []string
	
	limit struct {
		isUse  bool
		offset int64
		size   int64
	}
	orderBy struct {
		isUse bool
		col   string
		mode  OrderByMode // ASC or DESC
	}
	having struct {
		isUse     bool
		condition string
	}
	groupBy struct {
		isUse bool
		col   string
	}
}

func (w whereQueryBuild) ToQuery() string {
	var query = make([]string, 0, 6)
	query = append(query, w.prefix)
	if len(w.condition) > 0 {
		query = append(query, "WHERE "+strings.Join(w.condition, " AND "))
	}
	if w.having.isUse {
		query = append(query, w.having.condition)
	}
	if w.groupBy.isUse {
		query = append(query, fmt.Sprintf("GROUP BY %s", w.groupBy.col))
	}
	if w.orderBy.isUse {
		query = append(query, fmt.Sprintf("ORDER BY %s %s", w.orderBy.col, w.orderBy.mode))
	}
	if w.limit.isUse {
		query = append(query, fmt.Sprintf("LIMIT %d,%d", w.limit.offset, w.limit.size))
	}
	q := strings.Join(query, " ")
	w.table.logQuery(q)
	return q
}

type WhereOption func(wb *whereQueryBuild)

func Limit(offset, size int64) WhereOption {
	return func(wb *whereQueryBuild) {
		wb.limit.isUse = true
		wb.limit.offset = offset
		wb.limit.size = size
	}
}
func Condition(conOpt ...ConditionOption) WhereOption {
	return func(wb *whereQueryBuild) {
		con := &conditionQuery{table: wb.table}
		for _, setter := range conOpt {
			setter(con)
		}
		query := con.ToQuery()
		if query == "" {
			return
		}
		wb.condition = append(wb.condition, query)
	}
}
func Having(colName string, conOpt ...ConditionOption) WhereOption {
	return func(wb *whereQueryBuild) {
		column := fmt.Sprintf("%v", colName)
		wb.table.colValid(column)
		con := &conditionQuery{table: wb.table}
		for _, setter := range conOpt {
			setter(con)
		}
		wb.having.isUse = true
		wb.having.condition = con.ToQuery()
	}
}
func GroupBy(colName string) WhereOption {
	return func(wb *whereQueryBuild) {
		column := fmt.Sprintf("%v", colName)
		wb.table.colValid(column)
		wb.groupBy.isUse = true
		wb.groupBy.col = column
	}
}

type OrderByMode string

const (
	ASC  OrderByMode = "ASC"
	DESC OrderByMode = "DESC"
)

func OrderBy(colName string, mode OrderByMode) WhereOption {
	return func(wb *whereQueryBuild) {
		column := fmt.Sprintf("%v", colName)
		wb.table.colValid(column)
		wb.orderBy.isUse = true
		wb.orderBy.col = column
		wb.orderBy.mode = mode
	}
}

func newWhereQuery(table tableQuery, prefix string, opts []WhereOption) toQuery {
	w := &whereQueryBuild{table: table, prefix: prefix}
	for _, setter := range opts {
		setter(w)
	}
	return w
}

package mquery

import (
	"fmt"
	"strings"
)

type whereBuilder struct {
	qb *queryBuilder

	and   []string
	or    []string
	limit struct {
		isUse  bool
		offset int
		count  int
	}
	order struct {
		isUse bool
		col   string
		mode  string // ASC or DESC
	}
	join struct {
		isUse   bool
		table   string
		keyRoot string
		keyJoin string
	}
}

func newWhere(qb *queryBuilder) WhereBuilder {
	return &whereBuilder{
		qb: qb,
	}
}

type WhereBuilder interface {
	And(key string, ops Operator, value interface{}) WhereBuilder
	Or(key string, ops Operator, value interface{}) WhereBuilder
	In(col string, value ...interface{}) WhereBuilder

	Limit(count int) WhereBuilder
	OrderByASC(col string) WhereBuilder
	OrderByDESC(col string) WhereBuilder
	Search(col string, value string) WhereBuilder
	Join(tableName, keyRoot, keyJoin string) WhereBuilder

	IToQuery
}

func (wb *whereBuilder) Join(tableName, keyRoot, keyJoin string) WhereBuilder {
	wb.join.isUse = true
	wb.join.table = tableName
	wb.join.keyRoot = keyRoot
	wb.join.keyJoin = keyJoin
	return wb
}
func (wb *whereBuilder) Search(col string, value string) WhereBuilder {
	wb.and = append(wb.and, col+" like "+" '%"+value+"%'")
	return wb
}

func (wb *whereBuilder) And(key string, ops Operator, value interface{}) WhereBuilder {
	wb.and = append(wb.and, toString(key, ops, value))
	return wb
}

func (wb *whereBuilder) Or(key string, ops Operator, value interface{}) WhereBuilder {
	wb.or = append(wb.or, toString(key, ops, value))
	return wb
}
func (wb *whereBuilder) In(col string, value ...interface{}) WhereBuilder {
	if len(value) > 2 {
		for _, v := range value {
			if _, ok := v.(WhereBuilder); ok {
				panic("2 in")
			}
		}
	}

	listValue := []string{}
	for _, v := range value {
		listValue = append(listValue, interfaceToString(v))
	}
	wb.and = append(wb.and, fmt.Sprintf("%s IN (%s)", col, strings.Join(listValue, ",")))
	return wb
}

func (wb *whereBuilder) Limit(count int) WhereBuilder {
	wb.limit.isUse = true
	wb.limit.count = count
	return wb
}
func (wb *whereBuilder) Offset(off int) WhereBuilder {
	wb.limit.offset = off
	return wb
}

func (wb *whereBuilder) OrderByASC(col string) WhereBuilder {
	wb.order.isUse = true
	wb.order.mode = "ASC"
	wb.order.col = col
	return wb
}

func (wb *whereBuilder) OrderByDESC(col string) WhereBuilder {
	wb.order.isUse = true
	wb.order.mode = "DESC"
	wb.order.col = col
	return wb
}

func (wb *whereBuilder) ToQuery() string {
	query := []string{}
	if wb.join.isUse {
		query = append(query, fmt.Sprintf("JOIN %s ON %s.%s = %s.%s", wb.join.table, wb.qb.tableName, wb.join.keyRoot, wb.join.table, wb.join.keyJoin))
	}

	if len(wb.and) > 0 || len(wb.or) > 0 {
		qw := ""
		if and := strings.Join(wb.and, " and "); and != "" {
			qw += and + " "
		}
		if or := strings.Join(wb.or, " or "); or != "" {
			if qw != "" {
				qw += "or " + or
			} else {
				qw = or
			}
		}
		query = append(query, fmt.Sprintf("WHERE %s", qw))
	}
	if wb.order.isUse {
		query = append(query, fmt.Sprintf("ORDER BY %s %s", wb.order.col, wb.order.mode))
	}
	if wb.limit.isUse {
		query = append(query, fmt.Sprintf("LIMIT %d,%d", wb.limit.offset, wb.limit.count))
	}
	return strings.Join(query, " ")
}

package mquery

import (
	"fmt"
	"strings"
)

type WhereBuilder interface {
	And(key string, ops Operator, value interface{}) WhereBuilder
	Or(key string, ops Operator, value interface{}) WhereBuilder
	In(col string, value ...interface{}) WhereBuilder

	Limit(count int) WhereBuilder
	OrderByASC(col string) WhereBuilder
	OrderByDESC(col string) WhereBuilder
	Search(col string, value string) WhereBuilder

	IToQuery
}

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
}

func newWhereBuidler(qb *queryBuilder) WhereBuilder {
	return &whereBuilder{
		qb: qb,
	}
}

func (wb *whereBuilder) Search(col string, value string) WhereBuilder {
	wb.qb.colValid(col)
	wb.and = append(wb.and, col+" like "+" '%"+value+"%'")
	return wb
}

func (wb *whereBuilder) And(col string, ops Operator, value interface{}) WhereBuilder {
	wb.qb.colValid(col)
	wb.and = append(wb.and, toString(col, ops, value))
	return wb
}

func (wb *whereBuilder) Or(col string, ops Operator, value interface{}) WhereBuilder {
	wb.qb.colValid(col)
	wb.or = append(wb.or, toString(col, ops, value))
	return wb
}
func (wb *whereBuilder) In(col string, value ...interface{}) WhereBuilder {
	wb.qb.colValid(col)
	listValue := []string{}

	for _, v := range value {
		if _, ok := v.(WhereBuilder); ok {
			panic("dont use where build in here")
		}
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
	wb.qb.colValid(col)
	wb.order.isUse = true
	wb.order.mode = "ASC"
	wb.order.col = col
	return wb
}

func (wb *whereBuilder) OrderByDESC(col string) WhereBuilder {
	wb.qb.colValid(col)
	wb.order.isUse = true
	wb.order.mode = "DESC"
	wb.order.col = col
	return wb
}

func (wb *whereBuilder) ToQuery() string {
	query := []string{}
	if len(wb.and) > 0 || len(wb.or) > 0 {
		qw := []string{}
		if and := strings.Join(wb.and, " and "); and != "" {
			qw = append(qw, and)
		}
		if or := strings.Join(wb.or, " or "); or != "" {
			if len(qw) > 0 {
				or = "or " + or
			}
			qw = append(qw, or)
		}
		query = append(query, fmt.Sprintf("WHERE %s", strings.Join(qw, " ")))
	}
	if wb.order.isUse {
		query = append(query, fmt.Sprintf("ORDER BY %s %s", wb.order.col, wb.order.mode))
	}
	if wb.limit.isUse {
		query = append(query, fmt.Sprintf("LIMIT %d,%d", wb.limit.offset, wb.limit.count))
	}
	return strings.Join(query, " ")
}

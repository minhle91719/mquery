package mquery

import (
	"fmt"
	"strings"
)

type WhereBuilder interface {
	Condition(condition Condition) WhereBuilder
	Limit() WhereBuilder
	OrderByASC(col string) WhereBuilder
	OrderByDESC(col string) WhereBuilder
	// OrderByCase(listPriority ...string) WhereBuilder // uu tien giam dan
	Having(condition Condition) WhereBuilder
	GroupBy(col string) WhereBuilder
	IToQuery
}

type whereBuilder struct {
	qb *queryBuilder
	
	condition string
	limit     bool
	orderBy   struct {
		isUse bool
		col   string
		mode  string // ASC or DESC
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

func (wb *whereBuilder) Condition(condition Condition) WhereBuilder {
	wb.condition = condition.ToCondititonString()
	return wb
}

func (wb *whereBuilder) Having(condition Condition) WhereBuilder {
	wb.having.isUse = true
	wb.having.condition = condition.ToCondititonString()
	return wb
}

func (wb *whereBuilder) GroupBy(col string) WhereBuilder {
	wb.groupBy.isUse = true
	wb.groupBy.col = col
	return wb
}

func newWhereBuidler(qb *queryBuilder) WhereBuilder {
	return &whereBuilder{
		qb: qb,
	}
}

func (wb *whereBuilder) Limit() WhereBuilder {
	wb.limit = true
	return wb
}

func (wb *whereBuilder) OrderByASC(col string) WhereBuilder {
	wb.qb.colValid(col)
	wb.orderBy.isUse = true
	wb.orderBy.mode = "ASC"
	wb.orderBy.col = col
	return wb
}

func (wb *whereBuilder) OrderByDESC(col string) WhereBuilder {
	wb.qb.colValid(col)
	wb.orderBy.isUse = true
	wb.orderBy.mode = "DESC"
	wb.orderBy.col = col
	return wb
}

func (wb *whereBuilder) ToQuery() string {
	query := []string{}
	if wb.condition != "" {
		query = append(query, "WHERE "+wb.condition)
	}
	if wb.having.isUse {
		query = append(query, wb.having.condition)
	}
	if wb.groupBy.isUse {
		query = append(query, fmt.Sprintf("GROUP BY %s", wb.groupBy.col))
	}
	if wb.orderBy.isUse {
		query = append(query, fmt.Sprintf("ORDER BY %s %s", wb.orderBy.col, wb.orderBy.mode))
	}
	if wb.limit {
		query = append(query, "LIMIT ?,?")
	}
	return strings.Join(query, " ")
}

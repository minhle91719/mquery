package database

import (
	"fmt"
	"strings"
)

type QueryBuilder interface {
	// database.Where
	And(key string, value interface{}) QueryBuilder
	Or(key string, value interface{}) QueryBuilder
	Limit(count int) QueryBuilder
	OrderByASC(col string) QueryBuilder
	OrderByDESC(col string) QueryBuilder
	Search(col string, value string) QueryBuilder

	ToQuery() string // generate query
	// TODO: join
}

func Where(selectAll string) QueryBuilder {
	return &queryBuilder{
		selectStmt: selectAll,
	}
}

type queryBuilder struct {
	selectStmt string
	and        []string
	or         []string
	limit      struct {
		isUse  bool
		offset int
		count  int
	}
	order struct {
		isUse bool
		col   string
		isASC bool
	}
}

func (qb *queryBuilder) Search(col string, value string) QueryBuilder {
	qb.and = append(qb.and, col+" like "+" '%"+value+"%'")
	return qb
}

func (qb *queryBuilder) And(key string, value interface{}) QueryBuilder {
	qb.and = append(qb.and, toString(key, value))
	return qb
}

func (qb *queryBuilder) Or(key string, value interface{}) QueryBuilder {
	qb.or = append(qb.or, toString(key, value))
	return qb
}

func (qb *queryBuilder) Limit(count int) QueryBuilder {
	qb.limit.isUse = true
	qb.limit.count = count
	return qb
}
func (qb *queryBuilder) Offset(off int) QueryBuilder {
	qb.limit.offset = off
	return qb
}

func (qb *queryBuilder) OrderByASC(col string) QueryBuilder {
	qb.order.isUse = true
	qb.order.isASC = true
	qb.order.col = col
	return qb
}

func (qb *queryBuilder) OrderByDESC(col string) QueryBuilder {
	qb.order.isUse = true
	qb.order.isASC = false
	qb.order.col = col
	return qb
}
func (qb queryBuilder) ToQuery() string {
	delim := "or"
	if len(qb.and) == 0 || len(qb.or) == 0 {
		delim = ""
	}
	query := qb.selectStmt + fmt.Sprintf(" where %s %s %s", strings.Join(qb.and, " and "), delim, strings.Join(qb.or, " or "))
	if qb.order.isUse {
		ordMode := "DESC"
		if qb.order.isASC {
			ordMode = "ASC"
		}
		query += fmt.Sprintf(" order by %s %s", qb.order.col, ordMode)
	}
	if qb.limit.isUse {
		query += fmt.Sprintf(" limit %d,%d", qb.limit.offset, qb.limit.count)
	}
	return query
}
func toString(key string, value interface{}) string {
	query := fmt.Sprintf("%s", key)
	switch value.(type) {
	case int, uint:
		query += fmt.Sprintf(" = %d", value)
	case string:
		query += fmt.Sprintf(` = "%s"`, value)
	default:
		panic("unimplement ")
	}
	return query
}

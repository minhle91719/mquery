package database

import (
	"fmt"
	"strings"
)

type SelectQueryBuilder interface {
	// database.Where
	And(key string, value interface{}) SelectQueryBuilder
	Or(key string, value interface{}) SelectQueryBuilder
	Limit(count int) SelectQueryBuilder
	OrderByASC(col string) SelectQueryBuilder
	OrderByDESC(col string) SelectQueryBuilder
	Search(col string, value string) SelectQueryBuilder

	ToQuery() string // generate query
	// TODO: join
}

func Where(selectAll string) SelectQueryBuilder {
	return &queryBuidler{
		selectStmt: selectAll,
	}
}

type queryBuidler struct {
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
		mode  string // ASC or DESC
	}
}

func (qb *queryBuidler) Search(col string, value string) SelectQueryBuilder {
	qb.and = append(qb.and, col+" like "+" '%"+value+"%'")
	return qb
}

func (qb *queryBuidler) And(key string, value interface{}) SelectQueryBuilder {
	qb.and = append(qb.and, toString(key, value))
	return qb
}

func (qb *queryBuidler) Or(key string, value interface{}) SelectQueryBuilder {
	qb.or = append(qb.or, toString(key, value))
	return qb
}

func (qb *queryBuidler) Limit(count int) SelectQueryBuilder {
	qb.limit.isUse = true
	qb.limit.count = count
	return qb
}
func (qb *queryBuidler) Offset(off int) SelectQueryBuilder {
	qb.limit.offset = off
	return qb
}

func (qb *queryBuidler) OrderByASC(col string) SelectQueryBuilder {
	qb.order.isUse = true
	qb.order.mode = "ASC"
	qb.order.col = col
	return qb
}

func (qb *queryBuidler) OrderByDESC(col string) SelectQueryBuilder {
	qb.order.isUse = true
	qb.order.mode = "DESC"
	qb.order.col = col
	return qb
}
func (qb queryBuidler) ToQuery() string {
	delim := "or"
	if len(qb.and) == 0 || len(qb.or) == 0 {
		delim = ""
	}
	query := qb.selectStmt + fmt.Sprintf(" where %s %s %s", strings.Join(qb.and, " and "), delim, strings.Join(qb.or, " or "))
	if qb.order.isUse {
		query += fmt.Sprintf(" order by %s %s", qb.order.col, qb.order.mode)
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

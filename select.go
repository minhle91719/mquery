package mquery

import (
	"fmt"
	"strings"
)

type SelectQueryBuilder interface {
	//QueryBuilder
	Fields(col ...string) SelectQueryBuilder
	And(key string, value interface{}) SelectQueryBuilder
	Or(key string, value interface{}) SelectQueryBuilder
	Limit(count int) SelectQueryBuilder
	OrderByASC(col string) SelectQueryBuilder
	OrderByDESC(col string) SelectQueryBuilder
	Search(col string, value string) SelectQueryBuilder

	ToQuery() string // generate query
	// TODO: join
}

func Select(tableName string) SelectQueryBuilder {
	return &selectQueryBuidler{
		tableName: tableName,
	}
}

type selectQueryBuidler struct {
	tableName string
	fields    []string
	and       []string
	or        []string
	limit     struct {
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

func (qb *selectQueryBuidler) Fields(col ...string) SelectQueryBuilder {
	qb.fields = append(qb.fields, col...)
	return qb
}

func (qb *selectQueryBuidler) Search(col string, value string) SelectQueryBuilder {
	qb.and = append(qb.and, col+" like "+" '%"+value+"%'")
	return qb
}

func (qb *selectQueryBuidler) And(key string, value interface{}) SelectQueryBuilder {
	qb.and = append(qb.and, toString(key, value))
	return qb
}

func (qb *selectQueryBuidler) Or(key string, value interface{}) SelectQueryBuilder {
	qb.or = append(qb.or, toString(key, value))
	return qb
}

func (qb *selectQueryBuidler) Limit(count int) SelectQueryBuilder {
	qb.limit.isUse = true
	qb.limit.count = count
	return qb
}
func (qb *selectQueryBuidler) Offset(off int) SelectQueryBuilder {
	qb.limit.offset = off
	return qb
}

func (qb *selectQueryBuidler) OrderByASC(col string) SelectQueryBuilder {
	qb.order.isUse = true
	qb.order.mode = "ASC"
	qb.order.col = col
	return qb
}

func (qb *selectQueryBuidler) OrderByDESC(col string) SelectQueryBuilder {
	qb.order.isUse = true
	qb.order.mode = "DESC"
	qb.order.col = col
	return qb
}
func (qb selectQueryBuidler) ToQuery() string {
	var (
		where = ""
		field = ""
	)
	if len(qb.and) > 0 || len(qb.or) > 0 {
		qw := ""
		if and := strings.Join(qb.and, " and "); and != "" {
			qw += and + " "
		}
		if or := strings.Join(qb.or, " or "); or != "" {
			if qw != "" {
				qw += "or " + or
			} else {
				qw = or
			}
		}
		where = fmt.Sprintf("WHERE %s", qw)
	}
	if len(qb.fields) == 0 {
		field = "*"
	} else {
		field = strings.Join(qb.fields, ",")
	}

	query := fmt.Sprintf("SELECT %s FROM %s %s", field, qb.tableName, where)
	if qb.order.isUse {
		query = strings.TrimRight(query, " ")
		query += fmt.Sprintf(" ORDER BY %s %s", qb.order.col, qb.order.mode)
	}
	if qb.limit.isUse {
		query = strings.TrimRight(query, " ")
		query += fmt.Sprintf(" LIMIT %d,%d", qb.limit.offset, qb.limit.count)
	}
	// TODO: Join
	query = strings.TrimRight(query, " ")
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

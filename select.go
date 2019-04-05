package mquery

import (
	"fmt"
	"strings"
	"time"
)

type SelectQueryBuilder interface {
	//QueryBuilder
	Fields(col ...string) SelectQueryBuilder
	And(key string, value interface{}) SelectQueryBuilder
	Or(key string, value interface{}) SelectQueryBuilder
	In(col string, value ...interface{}) SelectQueryBuilder

	Limit(count int) SelectQueryBuilder
	OrderByASC(col string) SelectQueryBuilder
	OrderByDESC(col string) SelectQueryBuilder
	Search(col string, value string) SelectQueryBuilder
	Join(tableName, keyRoot, keyJoin string) SelectQueryBuilder

	ToQuery() string // generate query
	// TODO: join
}

func newSelectBuilder(qBuilder *queryBuilder) SelectQueryBuilder {
	return &selectQueryBuidler{
		qb: qBuilder,
	}
}

type selectQueryBuidler struct {
	qb     *queryBuilder
	fields []string
	and    []string
	or     []string
	limit  struct {
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

func (sqb *selectQueryBuidler) Fields(col ...string) SelectQueryBuilder {
	for _, v := range col {
		if sqb.qb.colValid(v) {
			sqb.fields = append(sqb.fields, col...)
		} else {
			panic("column " + v + " not exist . Please check " + sqb.qb.tableName + " QueryBuilder")
		}
	}

	return sqb
}
func (sqb *selectQueryBuidler) Join(tableName, keyRoot, keyJoin string) SelectQueryBuilder {
	sqb.join.isUse = true
	sqb.join.table = tableName
	sqb.join.keyRoot = keyRoot
	sqb.join.keyJoin = keyJoin
	return sqb

}

func (sqb *selectQueryBuidler) Search(col string, value string) SelectQueryBuilder {
	sqb.and = append(sqb.and, col+" like "+" '%"+value+"%'")
	return sqb
}

func (sqb *selectQueryBuidler) And(key string, value interface{}) SelectQueryBuilder {
	sqb.and = append(sqb.and, toString(key, value))
	return sqb
}

func (sqb *selectQueryBuidler) Or(key string, value interface{}) SelectQueryBuilder {
	sqb.or = append(sqb.or, toString(key, value))
	return sqb
}
func (sqb *selectQueryBuidler) In(col string, value ...interface{}) SelectQueryBuilder {
	listValue := []string{}
	for _, v := range value {
		listValue = append(listValue, interfaceToString(v))
	}
	sqb.and = append(sqb.and, fmt.Sprintf("%s IN (%s)", col, strings.Join(listValue, ",")))
	return sqb
}

func (sqb *selectQueryBuidler) Limit(count int) SelectQueryBuilder {
	sqb.limit.isUse = true
	sqb.limit.count = count
	return sqb
}
func (sqb *selectQueryBuidler) Offset(off int) SelectQueryBuilder {
	sqb.limit.offset = off
	return sqb
}

func (sqb *selectQueryBuidler) OrderByASC(col string) SelectQueryBuilder {
	sqb.order.isUse = true
	sqb.order.mode = "ASC"
	sqb.order.col = col
	return sqb
}

func (sqb *selectQueryBuidler) OrderByDESC(col string) SelectQueryBuilder {
	sqb.order.isUse = true
	sqb.order.mode = "DESC"
	sqb.order.col = col
	return sqb
}
func (sqb selectQueryBuidler) ToQuery() string {
	var (
		query = ""
		field = ""
	)
	if len(sqb.fields) == 0 {
		field = "*"
	} else {
		field = strings.Join(sqb.fields, ",")
	}
	query = fmt.Sprintf("SELECT %s FROM %s", field, sqb.qb.tableName)

	if sqb.join.isUse {
		query += " " + fmt.Sprintf("JOIN %s ON %s.%s = %s.%s", sqb.join.table, sqb.qb.tableName, sqb.join.keyRoot, sqb.join.table, sqb.join.keyJoin)
	}

	if len(sqb.and) > 0 || len(sqb.or) > 0 {
		qw := ""
		if and := strings.Join(sqb.and, " and "); and != "" {
			qw += and + " "
		}
		if or := strings.Join(sqb.or, " or "); or != "" {
			if qw != "" {
				qw += "or " + or
			} else {
				qw = or
			}
		}
		query += " " + fmt.Sprintf("WHERE %s", qw)
	}

	if sqb.order.isUse {
		query += " " + fmt.Sprintf("ORDER BY %s %s", sqb.order.col, sqb.order.mode)
	}
	if sqb.limit.isUse {
		query = " " + fmt.Sprintf("LIMIT %d,%d", sqb.limit.offset, sqb.limit.count)
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
func interfaceToString(value interface{}) string {
	switch value.(type) {
	case int, uint:
		return fmt.Sprintf("%d", value)
	case string:
		return fmt.Sprintf("%s", value)
	case time.Time:
		return value.(time.Time).String()
	default:
		panic("unimplement")
	}
}

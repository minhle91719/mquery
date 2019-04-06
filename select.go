package mquery

import (
	"fmt"
	"strings"
)

type SelectQueryBuilder interface {
	//QueryBuilder
	Fields(col ...string) SelectQueryBuilder
	Where(wb WhereBuilder) IToQuery
	IToQuery
}

func newSelectBuilder(qBuilder *queryBuilder) SelectQueryBuilder {
	return &selectQueryBuidler{
		qb: qBuilder,
	}
}

type selectQueryBuidler struct {
	qb     *queryBuilder
	fields []string

	where string
}

func (sqb *selectQueryBuidler) Fields(col ...string) SelectQueryBuilder {
	for _, v := range col {
		if sqb.qb.colValid(v) {
			sqb.fields = append(sqb.fields, v)
		} else {
			panic("column " + v + " not exist . Please check " + sqb.qb.tableName + " QueryBuilder")
		}
	}

	return sqb
}
func (sqb *selectQueryBuidler) Where(wb WhereBuilder) IToQuery {
	sqb.where = wb.ToQuery()
	return sqb
}

func (sqb *selectQueryBuidler) ToQuery() string {
	var (
		query = ""
		field = ""
	)
	if len(sqb.fields) == 0 {
		field = "*"
	} else {
		field = strings.Join(sqb.fields, ",")
	}
	query = fmt.Sprintf("SELECT %s FROM %s", field, sqb.qb.tableName) + " " + sqb.where
	return strings.TrimRight(query, " ")
}

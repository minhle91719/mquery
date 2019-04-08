package mquery

import (
	"fmt"
	"strings"
)

type SelectQueryBuilder interface {
	//QueryBuilder
	Fields(col ...string) SelectQueryBuilder
	//Join(tableName, keyRoot, keyJoin string) SelectQueryBuilder
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
	// join   struct {
	// 	isUse   bool
	// 	table   string
	// 	keyRoot string
	// 	keyJoin string
	// }
	where string
}

func (sqb *selectQueryBuidler) Fields(col ...string) SelectQueryBuilder {
	for _, v := range col {
		sqb.qb.colValid(v)
		sqb.fields = append(sqb.fields, v)
	}

	return sqb
}

// func (sqb *selectQueryBuidler) Join(tableName, keyRoot, keyJoin string) SelectQueryBuilder {
// 	sqb.qb.colValid(keyRoot)
// 	sqb.join.isUse = true
// 	sqb.join.table = tableName
// 	sqb.join.keyRoot = keyRoot
// 	sqb.join.keyJoin = keyJoin
// 	return sqb
// }
func (sqb *selectQueryBuidler) Where(wb WhereBuilder) IToQuery {
	sqb.where = wb.ToQuery()
	return sqb
}

func (sqb *selectQueryBuidler) ToQuery() string {
	var (
		query = []string{}
		field = ""
	)
	if len(sqb.fields) == 0 {
		field = "*"
	} else {
		field = strings.Join(sqb.fields, ",")
	}
	query = append(query, fmt.Sprintf("SELECT %s FROM %s", field, sqb.qb.tableName))
	// if sqb.join.isUse {
	// 	query = append(query, fmt.Sprintf("JOIN %s ON %s.%s = %s.%s", sqb.join.table, sqb.qb.tableName, sqb.join.keyRoot, sqb.join.table, sqb.join.keyJoin))
	// }
	if sqb.where != "" {
		query = append(query, sqb.where)
	}
	return strings.Join(query, " ")
}

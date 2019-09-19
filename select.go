package mquery

import (
	"fmt"
	"strings"
)

type SelectQueryBuilder interface {
	//QueryBuilder
	Fields(col ...interface{}) SelectQueryBuilder
	NotCheckFieldValid() SelectQueryBuilder
	//Join(tableName, keyRoot, keyJoin string) SelectQueryBuilder
	//CountWithDistict(colName, asName string) SelectQueryBuilder
	
	Where(wb WhereBuilder) IToQuery
	
	IToQuery
}

func newSelectBuilder(qBuilder *queryBuilder) SelectQueryBuilder {
	return &selectQueryBuidler{
		qb:           qBuilder,
		isCheckField: true,
	}
}

type selectQueryBuidler struct {
	qb           *queryBuilder
	isCheckField bool
	fields       []string
	count        struct {
		isUse   bool
		colName string
		asName  string
	}
	distinct struct {
		isUse   bool
		colName string
	}
	where string
}

func (sqb *selectQueryBuidler) NotCheckFieldValid() SelectQueryBuilder {
	sqb.isCheckField = false
	return sqb
}

func (sqb *selectQueryBuidler) Fields(col ...interface{}) SelectQueryBuilder {
	if len(col) == 1 && col[0] == "*" {
		sqb.fields = append(sqb.fields, sqb.qb.col...)
		return sqb
	}
	for _, v := range col {
		if sqb.isCheckField {
			sqb.qb.colValid(fmt.Sprint(v))
		}
		sqb.fields = append(sqb.fields, fmt.Sprint(v))
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
		query = make([]string, 0, len(sqb.fields)+2)
		field = ""
	)
	if len(sqb.fields) == 0 {
		return ""
	} else {
		field = strings.Join(sqb.fields, ",")
	}
	query = append(query, fmt.Sprintf("SELECT %s FROM %s", field, sqb.qb.tableName))
	if sqb.where != "" {
		query = append(query, sqb.where)
	}
	return strings.TrimSpace(strings.Join(query, " "))
}

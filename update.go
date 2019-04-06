package mquery

import (
	"fmt"
	"strings"
)

type UpdateQueryBuilder interface {
	Value(map[string]interface{}) UpdateQueryBuilder
	Where(wb WhereBuilder) IToQuery
	IToQuery
}
type updateQueryBuilder struct {
	qb       *queryBuilder
	mapValue map[string]interface{}
	where    string // TODO: using WHERE Select
}

func newUpdateBuilder(qb *queryBuilder) UpdateQueryBuilder {
	return &updateQueryBuilder{
		qb: qb,
	}
}
func (uqb *updateQueryBuilder) Value(mapValue map[string]interface{}) UpdateQueryBuilder {
	for k := range mapValue {
		uqb.qb.colValid(k)
	}
	uqb.mapValue = mapValue
	return uqb
}
func (uqb *updateQueryBuilder) Where(wb WhereBuilder) IToQuery {
	uqb.where = wb.ToQuery()
	return uqb
}
func (uqb *updateQueryBuilder) ToQuery() string {
	query := ""
	listUpdate := []string{}
	for k, v := range uqb.mapValue {
		uqb.qb.colValid(k)
		listUpdate = append(listUpdate, toString(k, Equal, v))
	}
	query = fmt.Sprintf("UPDATE %s SET %s %s", uqb.qb.tableName, strings.Join(listUpdate, ","), uqb.where)
	return strings.TrimRight(query, " ")
}

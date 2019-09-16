package mquery

import (
	"fmt"
	"strings"
)

type UpdateQueryBuilder interface {
	Values(mapValue map[string]interface{}) UpdateQueryBuilder
	Fields(listField ...string) UpdateQueryBuilder
	Where(wb WhereBuilder) IToQuery
	IToQuery
}
type updateQueryBuilder struct {
	qb      *queryBuilder
	listCol []string
	field   map[string]interface{}
	where   string // TODO: using WHERE Select
}

func (uqb *updateQueryBuilder) Values(mapValue map[string]interface{}) UpdateQueryBuilder {
	uqb.field = mapValue
	return uqb
}

func newUpdateBuilder(qb *queryBuilder) UpdateQueryBuilder {
	return &updateQueryBuilder{
		qb: qb,
	}
}
func (uqb *updateQueryBuilder) Fields(mapValue ...string) UpdateQueryBuilder {
	for _, v := range mapValue {
		uqb.qb.colValid(v)
	}
	uqb.listCol = mapValue
	return uqb
}
func (uqb *updateQueryBuilder) Where(wb WhereBuilder) IToQuery {
	uqb.where = wb.ToQuery()
	return uqb
}
func (uqb *updateQueryBuilder) ToQuery() string {
	var (
		query      []string
		listUpdate []string
	)
	for _, v := range uqb.listCol {
		listUpdate = append(listUpdate, toString(v, Equal, nil))
	}
	for k, v := range uqb.field {
		listUpdate = append(listUpdate, toString(k, Equal, v))
	}
	query = append(query, fmt.Sprintf("UPDATE %s SET %s", uqb.qb.tableName, strings.Join(listUpdate, ",")))
	if uqb.where != "" {
		query = append(query, uqb.where)
	}
	return strings.TrimSpace(strings.Join(query, " "))
}

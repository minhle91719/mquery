package mquery

import (
	"fmt"
	"html"
	"time"
)

type QueryBuilder interface {
	Fields(col map[string]bool) QueryBuilder // col name and not null
	InsertBuilder() InsertQueryBuilder
	SelectBuilder() SelectQueryBuilder
	WhereBuilder() WhereBuilder
	UpdateBuilder() UpdateQueryBuilder

	colValid(nameCol string)
}

type IToQuery interface {
	ToQuery() string
}

type Operator string

const (
	Equal    Operator = "="
	Less     Operator = "<"
	Greater  Operator = ">"
	NotEqual Operator = "<>"
)

type queryBuilder struct {
	tableName string
	col       map[string]bool
}

func NewTable(name string) QueryBuilder {
	qb := &queryBuilder{
		tableName: name,
	}
	return qb
}
func (qb *queryBuilder) Fields(mapCol map[string]bool) QueryBuilder {
	qb.col = mapCol
	return qb
}
func (qb *queryBuilder) InsertBuilder() InsertQueryBuilder {
	return newInsertBuilder(qb)
}
func (qb *queryBuilder) SelectBuilder() SelectQueryBuilder {
	return newSelectBuilder(qb)
}
func (qb *queryBuilder) WhereBuilder() WhereBuilder {
	return newWhereBuidler(qb)
}
func (qb *queryBuilder) UpdateBuilder() UpdateQueryBuilder {
	return newUpdateBuilder(qb)
}

func (qb queryBuilder) colValid(name string) {
	if _, ok := qb.col[name]; ok {
		return
	}
	panic("column " + name + " not exist . Please check " + qb.tableName + " QueryBuilder")
}
func toString(key string, ops Operator, value interface{}) string {
	return fmt.Sprintf("%s %s %s", key, ops, interfaceToString(value))
}
func interfaceToString(value interface{}) string {
	result := ""
	switch value.(type) {
	case int, uint, bool:
		result = fmt.Sprintf("%d", value)
	case string:
		result = fmt.Sprintf(`"%s"`, html.EscapeString(fmt.Sprintf("%s", value)))
	case time.Time:
		result = value.(time.Time).String()
	case SelectQueryBuilder, WhereBuilder:
		result = fmt.Sprintf("%s", value.(IToQuery).ToQuery())
	default:
		panic("unimplement")
	}
	return result
}

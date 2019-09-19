package mquery

import (
	"fmt"
	"html"
	"strings"
	"time"
)

type QueryBuilder interface {
	Fields(col []string) QueryBuilder // col name and not null
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
	Equal            Operator = "="
	Less             Operator = "<"
	GreaterThanEqual Operator = ">="
	LessThanEqual    Operator = "<="
	Greater          Operator = ">"
	NotEqual         Operator = "<>"
	
	Like Operator = "LIKE"
)

var (
	mapFuncSql = map[string]struct{}{
		"now()": struct{}{},
	}
)

type queryBuilder struct {
	tableName string
	col       []string
}

func NewTable(name string) QueryBuilder {
	qb := &queryBuilder{
		tableName: name,
	}
	return qb
}
func (qb *queryBuilder) Fields(mapCol []string) QueryBuilder {
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

func (qb *queryBuilder) colValid(name string) {
	name = replaceToken.Replace(name)
	for _, v := range qb.col {
		if v == name {
			return
		}
	}
	panic("column " + name + " not exist . Please check " + qb.tableName + " QueryBuilder")
}
func toString(key string, ops Operator, value interface{}) string {
	// check Key
	if _, ok := mapFuncSql[strings.ToLower(fmt.Sprint(value))]; ok {
		return fmt.Sprintf("%s %s %s", key, ops, strings.ToLower(fmt.Sprint(value)))
	}
	return fmt.Sprintf("%s %s %s", key, ops, interfaceToString(value))
}
func interfaceToString(value interface{}) string {
	result := ""
	switch value.(type) {
	case int, uint:
		result = fmt.Sprintf("%d", value)
	case string:
		result = fmt.Sprintf(`"%s"`, html.EscapeString(fmt.Sprintf("%s", value)))
	case time.Time:
		result = value.(time.Time).String()
	case SelectQueryBuilder, WhereBuilder:
		result = fmt.Sprintf("%s", value.(IToQuery).ToQuery())
	case bool:
		result = fmt.Sprint(value)
	case nil:
		result = "?"
	default:
		return fmt.Sprint(value)
	}
	return result
}

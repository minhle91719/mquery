package mquery

import (
	"fmt"
	"strings"
)

// TODO: binding map insert
type InsertQueryBuilder interface {
	Value(value ...string) IToQuery
}

type insertQueryBuilder struct {
	colValue []string
	
	qb *queryBuilder
}

func newInsertBuilder(qb *queryBuilder) InsertQueryBuilder {
	return &insertQueryBuilder{
		qb: qb,
	}
}
func (iqb *insertQueryBuilder) Value(value ...string) IToQuery {
	iqb.colValue = value
	return iqb
}
func (iqb *insertQueryBuilder) ToQuery() string {
	return fmt.Sprintf("INSERT INTO %s(%s) VALUE(%s)", iqb.qb.tableName, strings.Join(iqb.colValue, ","), genValueParam(len(iqb.colValue)))
}

func genValueParam(length int) (value string) {
	listValue := make([]string, 0, length)
	for i := 0; i < length; i++ {
		listValue = append(listValue, "?")
	}
	return strings.Join(listValue, ",")
}

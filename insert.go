package mquery

import (
	"fmt"
	"strings"
)

// TODO: binding map insert
type InsertQueryBuilder interface {
	Ignore() InsertQueryBuilder
	Value(value ...interface{}) IToQuery
}

type insertQueryBuilder struct {
	colValue []string
	ignore   bool
	values   bool
	qb       *queryBuilder
}

func (iqb *insertQueryBuilder) Ignore() InsertQueryBuilder {
	iqb.ignore = true
	return iqb
}

func newInsertBuilder(qb *queryBuilder) InsertQueryBuilder {
	return &insertQueryBuilder{
		qb: qb,
	}
}
func (iqb *insertQueryBuilder) Value(value ...interface{}) IToQuery {
	for _, v := range value {
		iqb.colValue = append(iqb.colValue, fmt.Sprint(v))
	}
	return iqb
}
func (iqb *insertQueryBuilder) ToQuery() string {
	first := "INSERT INTO"
	if iqb.ignore {
		first = "INSERT IGNORE INTO"
	}
	return fmt.Sprintf("%s %s(%s) VALUE(%s)", first, iqb.qb.tableName, strings.Join(iqb.colValue, ","), genValueParam(len(iqb.colValue)))
}

func genValueParam(length int) (value string) {
	listValue := make([]string, 0, length)
	for i := 0; i < length; i++ {
		listValue = append(listValue, "?")
	}
	return strings.Join(listValue, ",")
}

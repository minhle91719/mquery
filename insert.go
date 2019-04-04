package mquery

import (
	"fmt"
	"strings"
)

type InsertQueryBuilder interface {
	Value(value map[string]interface{}) InsertQueryBuilder

	ToQuery() string
}

type insertQueryBuilder struct {
	tableName string
	col       []string
	value     []interface{}
}

func NewInsertBuilder(tableName string) InsertQueryBuilder {
	return &insertQueryBuilder{
		tableName: tableName,
	}
}
func (iqb *insertQueryBuilder) Value(mapValue map[string]interface{}) InsertQueryBuilder {
	for k, v := range mapValue {
		iqb.col = append(iqb.col, k)
		iqb.value = append(iqb.value, v)
	}
	return iqb
}
func (iqb *insertQueryBuilder) ToQuery() string {
	listValue := []string{}
	for _, v := range iqb.value {
		listValue = append(listValue, interfaceToString(v))
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", iqb.tableName, strings.Join(iqb.col, ","), strings.Join(listValue, ","))
}

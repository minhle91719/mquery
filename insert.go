package mquery

import (
	"fmt"
	"strings"
)

// TODO: binding map insert
type InsertQueryBuilder interface {
	// example:
	// map de kiem tra tinh kha dung cua query sau chay check toan bo
	Column(col ...string) InsertQueryBuilder
	Value(value map[string]interface{}) InsertQueryBuilder

	IToQuery
}

type insertQueryBuilder struct {
	mapColValue map[string]interface{}

	qb *queryBuilder
}

func newInsertBuilder(qb *queryBuilder) InsertQueryBuilder {
	return &insertQueryBuilder{
		qb: qb,
	}
}
func (iqb *insertQueryBuilder) Column(col ...string) InsertQueryBuilder {
	for _, v := range col {
		iqb.mapColValue[v] = 0 // init
	}
	return iqb
}

func (iqb *insertQueryBuilder) Value(mapValue map[string]interface{}) InsertQueryBuilder {
	for k, v := range mapValue {
		if !iqb.qb.colValid(k) {
			panic("column not exist. Please check " + iqb.qb.tableName + " QueryBuilder")
		}

		key := strings.TrimLeft(k, ":")
		if value, ok := iqb.mapColValue[key]; ok && value != nil {
			iqb.mapColValue[key] = v
		} else {
			panic("cant find column : " + k)
		}
	}
	return iqb
}
func (iqb *insertQueryBuilder) ToQuery() string {
	listCol := []string{}
	listValue := []string{}
	for k, v := range iqb.mapColValue {
		listCol = append(listCol, k)
		listValue = append(listValue, interfaceToString(v))
	}
	return fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", iqb.qb.tableName, strings.Join(listCol, ","), strings.Join(listValue, ","))
}

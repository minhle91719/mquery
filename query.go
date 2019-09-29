package mquery

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type tableQuery struct {
	column    []string
	tableName string
	isLogger  bool
	logger    *logrus.Entry
}

func (tb tableQuery) Delete(opts ...WhereOption) toQuery {
	return newWhereQuery(tb, fmt.Sprintf("DELETE FROM %s", tb.tableName), opts)
}

func (tb tableQuery) Select(opts ...SelectOption) WhereQuery {
	return newSelect(tb, opts)
}

func (tb tableQuery) Update(opts ...UpdateOption) WhereQuery {
	return newUpdateQuery(tb, opts)
}

func (tb tableQuery) colValid(name interface{}) {
	col := fmt.Sprintf("%v", name)
	col = replaceToken.Replace(col)
	for _, v := range tb.column {
		if v == col {
			return
		}
	}
	logrus.Fatal("column " + col + " not exist . Please check " + tb.tableName + " QueryBuilder")
}
func (tb tableQuery) Insert(options ...InsertOption) toQuery {
	return newInsert(tb, options)
}
func (tb tableQuery) logQuery(q string) {
	if tb.isLogger {
		tb.logger.WithFields(logrus.Fields{
			"table": tb.tableName,
			"query": q,
		}).Info("mquery")
	}
}

type TableOption func(q *tableQuery)

func Column(value ...interface{}) TableOption {
	return func(q *tableQuery) {
		for _, v := range value {
			q.column = append(q.column, fmt.Sprintf("%v", v))
		}
	}
}

func WithLogger(l *logrus.Entry) TableOption {
	return func(q *tableQuery) {
		if l != nil {
			q.isLogger = true
			q.logger = l.WithFields(logrus.Fields{
				"infra": "mysql",
				"table": q.tableName,
			})
		}
	}
}

func NewQueryBuilder(tableName string, options ...TableOption) QueryBuild {
	q := &tableQuery{tableName: tableName}
	for _, setter := range options {
		setter(q)
	}
	return q
}

type QueryBuild interface {
	Insert(opts ...InsertOption) toQuery
	Select(opts ...SelectOption) WhereQuery
	Update(opts ...UpdateOption) WhereQuery
	Delete(opts ...WhereOption) toQuery
}

type WhereQuery interface {
	Where(...WhereOption) toQuery
}

type toQuery interface {
	ToQuery() string
}

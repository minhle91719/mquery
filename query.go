package mquery

type QueryBuilder interface {
	InsertBuilder() InsertQueryBuilder
	SelectBuilder() SelectQueryBuilder
	// UpdateBuilder()
}
type queryBuilder struct {
	tableName string
}

func NewTable(name string) QueryBuilder {
	return &queryBuilder{
		tableName: name,
	}
}

func (qb *queryBuilder) Table(name string) QueryBuilder {
	qb.tableName = name
	return qb
}
func (qb *queryBuilder) InsertBuilder() InsertQueryBuilder {
	return NewInsertBuilder(qb.tableName)
}
func (qb *queryBuilder) SelectBuilder() SelectQueryBuilder {
	return NewSelectBuilder(qb.tableName)
}

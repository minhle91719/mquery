package mquery

type QueryBuilder interface {
	Table(name string) QueryBuilder
	InsertBuilder() InsertQueryBuilder
	SelectBuilder() SelectQueryBuilder
	// UpdateBuilder()
}
type queryBuilder struct {
	tableName string
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

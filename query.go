package mquery

type QueryBuilder interface {
	Fields(col map[string]bool) QueryBuilder // col name and not null
	InsertBuilder() InsertQueryBuilder
	SelectBuilder() SelectQueryBuilder
	// UpdateBuilder()

	colValid(nameCol string) bool
}
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
func (qb queryBuilder) colValid(name string) bool {
	if _, ok := qb.col[name]; ok {
		return true
	}
	return false
}

package mquery

type UpdateQueryBuilder interface {
}
type updateQueryBuilder struct {
	qb          *queryBuilder
	valueUpdate []string
	where       string // TODO: using WHERE Select
}

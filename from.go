package mquery

import "fmt"

type fromQueryBuild struct {
	table         tableQuery
	forUpdate     bool
	prefix        string
	isUsingNested bool
	query         string
	as            string
}
type FromOption func(fo *fromQueryBuild)

func FromAs(query string, as interface{}) FromOption {
	return func(fo *fromQueryBuild) {
		fo.isUsingNested = true
		fo.query = query
		fo.as = fmt.Sprintf("%s", as)
	}
}

func newFromQuery(table tableQuery, isForUpdate bool, prefix string, options []FromOption) *fromQueryBuild {
	fo := &fromQueryBuild{
		prefix:    prefix,
		table:     table,
		forUpdate: isForUpdate,
	}
	for _, v := range options {
		v(fo)
	}
	return fo
}

func (fo *fromQueryBuild) ToQuery() string {
	fromResource := fo.table.tableName
	if fo.isUsingNested {
		fromResource = fmt.Sprintf("(%s) as %s", fo.query, fo.as)
	}
	return fmt.Sprintf("%s FROM %s", fo.prefix, fromResource)
}
func (fo *fromQueryBuild) Where(opts ...WhereOption) toQuery {
	if fo.forUpdate {
		opts = append(opts, forUpdate())
	}
	return newWhereQuery(fo.table, fo.ToQuery(), opts)
}

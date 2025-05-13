package textsearch

import "github.com/futugyou/yomawari/semantic_kernel/abstractions/vectordata"

type TextSearchFilter struct {
	filterClauses []vectordata.FilterClause
}

func (t TextSearchFilter) GetFilterClauses() []vectordata.FilterClause {
	return t.filterClauses
}

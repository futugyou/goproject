package abstractions

type FilterClause interface {
	IsFilterClause()
}

type AnyTagEqualToFilterClause struct {
	FieldName string
	Value     string
}

func (f AnyTagEqualToFilterClause) IsFilterClause() {}

type EqualToFilterClause struct {
	FieldName string
	Value     any
}

func (f EqualToFilterClause) IsFilterClause() {}

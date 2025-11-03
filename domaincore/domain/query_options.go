package domain

type QueryOptions struct {
	OrderBy map[string]int
	Page    int
	Limit   int
	Filter  FilterExpr
}

func NewQueryOptions(page *int, limit *int, orderBy map[string]int, filter FilterExpr) *QueryOptions {
	opts := &QueryOptions{
		OrderBy: orderBy,
		Page:    1,
		Limit:   100,
		Filter:  filter,
	}

	if page != nil && *page > 1 {
		opts.Page = *page
	}

	if limit != nil && *limit > 1 && *limit < 100 {
		opts.Limit = *limit
	}

	return opts
}

type FilterExpr any

type And []FilterExpr
type Or []FilterExpr
type Not struct {
	Expr FilterExpr
}

type Eq struct {
	Field string
	Value any
}

type Ne struct {
	Field string
	Value any
}

type Gt struct {
	Field string
	Value any
}

type Gte struct {
	Field string
	Value any
}

type Lt struct {
	Field string
	Value any
}

type Lte struct {
	Field string
	Value any
}

type In struct {
	Field  string
	Values []any
}

type Nin struct {
	Field  string
	Values []any
}

type Like struct {
	Field           string
	Pattern         string
	CaseInsensitive bool
}

type QueryBuilder struct {
	expr FilterExpr
}

func NewQuery() *QueryBuilder {
	return &QueryBuilder{}
}

func (b *QueryBuilder) Eq(field string, value any) *QueryBuilder {
	b.add(Eq{Field: field, Value: value})
	return b
}

func (b *QueryBuilder) Ne(field string, value any) *QueryBuilder {
	b.add(Ne{Field: field, Value: value})
	return b
}

func (b *QueryBuilder) Gt(field string, value any) *QueryBuilder {
	b.add(Gt{Field: field, Value: value})
	return b
}

func (b *QueryBuilder) Gte(field string, value any) *QueryBuilder {
	b.add(Gte{Field: field, Value: value})
	return b
}

func (b *QueryBuilder) Lt(field string, value any) *QueryBuilder {
	b.add(Lt{Field: field, Value: value})
	return b
}

func (b *QueryBuilder) Lte(field string, value any) *QueryBuilder {
	b.add(Lte{Field: field, Value: value})
	return b
}

func (b *QueryBuilder) In(field string, values ...any) *QueryBuilder {
	b.add(In{Field: field, Values: values})
	return b
}

func (b *QueryBuilder) Nin(field string, values ...any) *QueryBuilder {
	b.add(Nin{Field: field, Values: values})
	return b
}

func (b *QueryBuilder) Like(field, pattern string, caseInsensitive bool) *QueryBuilder {
	b.add(Like{Field: field, Pattern: pattern, CaseInsensitive: caseInsensitive})
	return b
}

func (b *QueryBuilder) And(subs ...FilterExpr) *QueryBuilder {
	b.add(And(subs))
	return b
}

func (b *QueryBuilder) Or(subs ...FilterExpr) *QueryBuilder {
	b.add(Or(subs))
	return b
}

func (b *QueryBuilder) Not(sub FilterExpr) *QueryBuilder {
	b.add(Not{Expr: sub})
	return b
}

func (b *QueryBuilder) add(f FilterExpr) {
	if b.expr == nil {
		b.expr = f
	} else {
		b.expr = And{b.expr, f}
	}
}

func (b *QueryBuilder) Build() FilterExpr {
	return b.expr
}

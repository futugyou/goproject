package domain

import (
	"context"
)

const DATA_NOT_FOUND_MESSAGE string = "data not found"

type Repository[Aggregate AggregateRoot] interface {
	FindByID(ctx context.Context, id string) (*Aggregate, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, aggregate Aggregate) error
	Insert(ctx context.Context, aggregate Aggregate) error
	Find(ctx context.Context, options *QueryOptions) ([]Aggregate, error)
}

// e.g. (id=="007" AND name=="Alice") OR (status=="active")
//
//	opts := &QueryOptions{
//	    LogicFilter: map[string]any{
//	        "$or": []any{
//	            map[string]any{
//	                "$and": []any{
//	                    map[string]any{"id": "007"},
//	                    map[string]any{"name": "Alice"},
//	                },
//	            },
//	            map[string]any{"status": "active"},
//	        },
//	    },
//	}
//
//	{
//		"$or": [
//		  { "$and": [ { "id": "007" }, { "name": "Alice" } ] },
//		  { "status": "active" }
//		]
//	}
type QueryOptions struct {
	OrderBy     map[string]int
	Page        int
	Limit       int
	Filters     map[string]any // only $and
	LogicFilter any            // Logical filters are used for complex query conditions, such as $or, $and, $not, and $nor. They have higher priority than regular `Filters`
}

func NewQueryOptions(page *int, limit *int, orderBy map[string]int, filters map[string]any) *QueryOptions {
	opts := &QueryOptions{
		OrderBy: orderBy,
		Page:    1,
		Limit:   100,
		Filters: filters,
	}

	if page != nil && *page > 1 {
		opts.Page = *page
	}

	if limit != nil && *limit > 1 && *limit < 100 {
		opts.Limit = *limit
	}

	return opts
}

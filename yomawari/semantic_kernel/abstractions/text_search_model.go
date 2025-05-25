package abstractions

import (
	"fmt"
)

type MapFromResultToString = func(result any) (*string, error)

type ITextSearchStringMapper interface {
	MapFromResultToString(result any) (*string, error)
}

type TextSearchStringMapper struct {
	mapFromResultToString MapFromResultToString
}

func NewTextSearchStringMapper(mapFromResultToString MapFromResultToString) *TextSearchStringMapper {
	return &TextSearchStringMapper{mapFromResultToString: mapFromResultToString}
}

func (t *TextSearchStringMapper) MapFromResultToString(result any) (*string, error) {
	if t == nil || t.mapFromResultToString == nil {
		return nil, fmt.Errorf("mapFromResultToString is nil")
	}

	return t.mapFromResultToString(result)
}

type TextSearchFilter struct {
	filterClauses []FilterClause
}

func (t TextSearchFilter) GetFilterClauses() []FilterClause {
	return t.filterClauses
}

const TextSearchOptionsDefaultTop int = 5

type TextSearchOptions struct {
	IncludeTotalCount bool
	Filter            *TextSearchFilter
	Top               int
	Skip              int
}

func NewTextSearchOptions() *TextSearchOptions {
	return &TextSearchOptions{
		IncludeTotalCount: false,
		Filter:            nil,
		Top:               TextSearchOptionsDefaultTop,
		Skip:              0,
	}
}

type TextSearchResult struct {
	Name  string
	Link  string
	Value string
}

func NewTextSearchResult(value string) *TextSearchResult {
	return &TextSearchResult{Value: value}
}

type KernelSearchResults[T any] struct {
	TotalCount int64
	Metadata   map[string]any
	Results    []T
}

func NewKernelSearchResults[T any](results []T, totalCount int64, metadata map[string]any) *KernelSearchResults[T] {
	return &KernelSearchResults[T]{
		TotalCount: totalCount,
		Metadata:   metadata,
		Results:    results,
	}
}

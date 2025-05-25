package abstractions

import (
	"fmt"
)

type MapFromResultToTextSearchResult func(result any) (*TextSearchResult, error)

type ITextSearchResultMapper interface {
	MapFromResultToTextSearchResult(result any) (*TextSearchResult, error)
}

type TextSearchResultMapper struct {
	mapFromResultToTextSearchResult MapFromResultToTextSearchResult
}

func NewTextSearchResultMapper(mapFromResultToTextSearchResult MapFromResultToTextSearchResult) *TextSearchResultMapper {
	return &TextSearchResultMapper{mapFromResultToTextSearchResult: mapFromResultToTextSearchResult}
}

func (m *TextSearchResultMapper) MapFromResultToTextSearchResult(result any) (*TextSearchResult, error) {
	if m == nil || m.mapFromResultToTextSearchResult == nil {
		return nil, fmt.Errorf("mapFromResultToTextSearchResult is nil")
	}

	return m.mapFromResultToTextSearchResult(result)
}

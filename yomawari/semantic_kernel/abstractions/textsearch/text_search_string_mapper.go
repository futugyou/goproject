package textsearch

import "fmt"

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

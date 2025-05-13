package textsearch

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

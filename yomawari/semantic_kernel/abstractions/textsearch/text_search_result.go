package textsearch

type TextSearchResult struct {
	Name  string
	Link  string
	Value string
}

func NewTextSearchResult(value string) *TextSearchResult {
	return &TextSearchResult{Value: value}
}

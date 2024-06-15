package extensions

// search condition
type Search struct {
	Sort   map[string]int
	Page   int
	Size   int
	Filter map[string]string
}

func NewSearch(page *int, size *int, sort map[string]int, filter map[string]string) *Search {
	s := &Search{
		Sort:   sort,
		Page:   1,
		Size:   100,
		Filter: filter,
	}
	if page != nil && *page > 1 {
		s.Page = *page
	}
	if size != nil && *size > 1 && *size < 100 {
		s.Size = *size
	}
	return s
}

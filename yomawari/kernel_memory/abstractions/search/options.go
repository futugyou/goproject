package search

type SearchOptions struct {
	Stream bool
}

func (s *SearchOptions) Clone() *SearchOptions {
	if s == nil {
		return nil
	}

	return &SearchOptions{
		Stream: s.Stream,
	}
}

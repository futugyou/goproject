package textsearch

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

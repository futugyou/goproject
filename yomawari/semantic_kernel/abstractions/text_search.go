package abstractions

import "context"

type ITextSearch interface {
	Search(ctx context.Context, query string, searchOptions TextSearchOptions) (*KernelSearchResults[string], error)
	GetTextSearchResults(ctx context.Context, query string, searchOptions TextSearchOptions) (*KernelSearchResults[TextSearchResult], error)
	GetSearchResults(ctx context.Context, query string, searchOptions TextSearchOptions) (*KernelSearchResults[any], error)
}

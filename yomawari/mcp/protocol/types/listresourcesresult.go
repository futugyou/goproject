package types

type ListResourcesResult struct {
	PaginatedResult
	Resources []Resource `json:"resources"`
}

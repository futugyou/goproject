package protocol

type ListResourcesResult struct {
	PaginatedResult `json:",inline"`
	Resources       []Resource `json:"resources"`
}

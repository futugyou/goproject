package types
 
type ListResourcesResult struct {
	PaginatedResult `json:",inline"`
	Resources                 []Resource `json:"resources"`
}

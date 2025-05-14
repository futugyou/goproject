package types
 
type ListToolsResult struct {
	 PaginatedResult `json:",inline"`
	Tools                     []Tool `json:"tools"`
}

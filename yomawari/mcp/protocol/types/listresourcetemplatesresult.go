package types
 
type ListResourceTemplatesResult struct {
	PaginatedResult `json:",inline"`
	ResourceTemplates         []ResourceTemplate `json:"resourceTemplates"`
}

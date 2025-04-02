package types

type ListResourceTemplatesResult struct {
	PaginatedResult
	ResourceTemplates []ResourceTemplate `json:"resourceTemplates"`
}

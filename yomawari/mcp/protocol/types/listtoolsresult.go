package types

type ListToolsResult struct{
	PaginatedResult
	Tools []Tool`json:"tools"`
}
package models

type SearchQuery struct {
	Index            *string                `json:"index"`
	Query            string                 `json:"query"`
	Filters          []MemoryFilter         `json:"filters"`
	MinRelevance     float64                `json:"minRelevance"`
	Limit            int64                  `json:"limit"`
	ContextArguments map[string]interface{} `json:"args"`
}

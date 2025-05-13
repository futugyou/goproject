package models

type MemoryQuery struct {
	Index            *string                `json:"index"`
	Question         string                 `json:"question"`
	Filters          []MemoryFilter         `json:"filters"`
	MinRelevance     float64                `json:"minRelevance"`
	Stream           bool                   `json:"stream"`
	ContextArguments map[string]interface{} `json:"args"`
}

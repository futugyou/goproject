package server

type ServerCapabilities struct {
	Experimental map[string]interface{} `json:"experimental"`
	Logging      *LoggingCapability     `json:"logging"`
	Prompts      *PromptsCapability     `json:"prompts"`
	Resources    *ResourcesCapability   `json:"resources"`
	Tools        *ToolsCapability       `json:"tools"`
}

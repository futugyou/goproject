package server

import "github.com/futugyou/yomawari/mcp/protocol"

type ServerCapabilities struct {
	Experimental         map[string]interface{}                  `json:"experimental"`
	Logging              *LoggingCapability                      `json:"logging"`
	Prompts              *PromptsCapability                      `json:"prompts"`
	Resources            *ResourcesCapability                    `json:"resources"`
	Tools                *ToolsCapability                        `json:"tools"`
	Completions          *CompletionsCapability                  `json:"completions"`
	NotificationHandlers map[string]protocol.NotificationHandler `json:"-"`
}

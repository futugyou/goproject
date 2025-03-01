package chatcompletion

type ChatMessage struct {
	Role                 ChatRole               `json:"role"`
	Message              string                 `json:"message"`
	Text                 *string                `json:"-"`
	Contents             []interface{}          `json:"contents"`
	AuthorName           *string                `json:"authorName"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

type ChatRole string

const (
	RoleSystem    ChatRole = "system"
	RoleAssistant ChatRole = "assistant"
	RoleUser      ChatRole = "user"
	RoleTool      ChatRole = "tool"
)

package chatcompletion

type ChatMessage struct {
	// system assistant user tool
	Role ChatRole `json:"role"`
	// The text of the message.
	Message string `json:"message"`
}

type ChatRole string

const (
	RoleSystem    ChatRole = "system"
	RoleAssistant ChatRole = "assistant"
	RoleUser      ChatRole = "user"
	RoleTool      ChatRole = "tool"
)

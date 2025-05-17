package chat_completion

type AuthorRole string

const (
	AuthorRoleDeveloper AuthorRole = "developer"
	AuthorRoleSystem    AuthorRole = "system"
	AuthorRoleAssistant AuthorRole = "assistant"
	AuthorRoleUser      AuthorRole = "user"
	AuthorRoleTool      AuthorRole = "tool"
)

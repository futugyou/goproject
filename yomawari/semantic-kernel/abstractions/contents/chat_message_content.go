package contents

import "encoding/base64"

type AuthorRole string

const (
	AuthorNameDeveloper AuthorRole = "developer"
	AuthorRoleSystem    AuthorRole = "system"
	AuthorRoleAssistant AuthorRole = "assistant"
	AuthorRoleUser      AuthorRole = "user"
	AuthorRoleTool      AuthorRole = "tool"
)

type ChatMessageContentItemCollection struct {
	// TODO
}
type ChatMessageContent struct {
	MimeType   string                           `json:"mimeType"`
	ModelId    string                           `json:"modelId"`
	Metadata   map[string]any                   `json:"metadata"`
	AuthorName string                           `json:"authorName"`
	Role       AuthorRole                       `json:"role"`
	Items      ChatMessageContentItemCollection `json:"items"`
	Content    string                           `json:"-"`
	Encoding   *base64.Encoding                 `json:"-"`
	Source     any                              `json:"-"`
}

func (ChatMessageContent) Type() string {
	return "chatMessage"
}

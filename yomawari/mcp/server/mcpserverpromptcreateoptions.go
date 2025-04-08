package server

type McpServerPromptCreateOptions struct {
	Name        *string
	Description *string
}

func (m *McpServerPromptCreateOptions) Clone() *McpServerPromptCreateOptions {
	return &McpServerPromptCreateOptions{
		Name:        m.Name,
		Description: m.Description,
	}
}

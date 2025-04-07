package server

type McpServerPromptCreateOptions struct {
	Name       *string
	Descriptio *string
}

func (m *McpServerPromptCreateOptions) Clone() *McpServerPromptCreateOptions {
	return &McpServerPromptCreateOptions{
		Name:       m.Name,
		Descriptio: m.Descriptio,
	}
}

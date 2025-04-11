package server

type McpServerToolCreateOptions struct {
	Name        *string
	Description *string
	Title       *string
	Destructive *bool
	Idempotent  *bool
	OpenWorld   *bool
	ReadOnly    *bool
}

func (m *McpServerToolCreateOptions) Clone() *McpServerToolCreateOptions {
	return &McpServerToolCreateOptions{
		Name:        m.Name,
		Description: m.Description,
		Title:       m.Title,
		Destructive: m.Destructive,
		Idempotent:  m.Idempotent,
		OpenWorld:   m.OpenWorld,
		ReadOnly:    m.ReadOnly,
	}
}

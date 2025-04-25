package abstractions

type HostedWebSearchTool struct {
	*BaseAITool
}

func NewHostedWebSearchTool() *HostedWebSearchTool {
	return &HostedWebSearchTool{
		BaseAITool: NewBaseAITool(),
	}
}

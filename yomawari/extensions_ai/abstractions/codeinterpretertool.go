package abstractions

type CodeInterpreterTool struct {
	*BaseAITool
}

func NewCodeInterpreterTool() *CodeInterpreterTool {
	return &CodeInterpreterTool{
		BaseAITool: NewBaseAITool(),
	}
}

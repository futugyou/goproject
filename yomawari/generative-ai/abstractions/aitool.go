package abstractions

// TODO: json schema
type AITool interface {
	GetName() string
	GetDescription() string
	GetParameters() map[string]interface{}
	GetAdditionalProperties() map[string]interface{}
}

type BaseAITool struct {
	Name                 string
	Description          string
	AdditionalProperties map[string]interface{}
}

func NewBaseAITool(name string) BaseAITool {
	return BaseAITool{
		Name:                 name,
		Description:          "",
		AdditionalProperties: make(map[string]interface{}),
	}
}

func (t BaseAITool) GetName() string {
	return t.Name
}

func (t BaseAITool) GetDescription() string {
	return t.Description
}

func (t BaseAITool) GetAdditionalProperties() map[string]interface{} {
	return t.AdditionalProperties
}

func (t BaseAITool) GetParameters() map[string]interface{} {
	return map[string]interface{}{}
}

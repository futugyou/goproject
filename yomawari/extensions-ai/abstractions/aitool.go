package abstractions

var _ AITool = (*BaseAITool)(nil)

type AITool interface {
	GetName() string
	GetDescription() string
	GetAdditionalProperties() map[string]interface{}
	//Obsolete
	GetParameters() map[string]interface{}
}

type BaseAITool struct {
	Name                 string
	Description          string
	AdditionalProperties map[string]interface{}
}

func NewBaseAITool(name string) *BaseAITool {
	return &BaseAITool{
		Name:                 name,
		Description:          "",
		AdditionalProperties: make(map[string]interface{}),
	}
}

func (t *BaseAITool) GetName() string {
	return t.Name
}

func (t *BaseAITool) GetDescription() string {
	return t.Description
}

func (t *BaseAITool) GetAdditionalProperties() map[string]interface{} {
	return t.AdditionalProperties
}

func (t *BaseAITool) GetParameters() map[string]interface{} {
	return map[string]interface{}{}
}

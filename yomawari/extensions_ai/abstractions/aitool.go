package abstractions

import "reflect"

var _ AITool = (*BaseAITool)(nil)

type AITool interface {
	GetName() string
	GetDescription() string
	GetAdditionalProperties() map[string]interface{}
}

type BaseAITool struct {
	Name                 string
	Description          string
	AdditionalProperties map[string]interface{}
}

func NewBaseAITool() *BaseAITool {
	return &BaseAITool{
		Description:          "",
		AdditionalProperties: make(map[string]interface{}),
	}
}

func (t *BaseAITool) GetName() string {
	return reflect.TypeOf(t).Name()
}

func (t *BaseAITool) GetDescription() string {
	return t.Description
}

func (t *BaseAITool) GetAdditionalProperties() map[string]interface{} {
	return t.AdditionalProperties
}

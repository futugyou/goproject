package functions

type KernelFunction interface {
	GetName() string
	GetPluginName() string
	GetAdditionalProperties() map[string]interface{}
	GetDescription() string
}

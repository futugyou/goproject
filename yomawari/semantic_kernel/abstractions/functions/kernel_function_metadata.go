package functions

type KernelFunctionMetadata struct {
	Name                 string
	PluginName           string
	Description          string
	Parameters           []KernelParameterMetadata
	ReturnParameter      KernelReturnParameterMetadata
	AdditionalProperties map[string]interface{}
}

func NewKernelFunctionMetadata(name string) *KernelFunctionMetadata {
	return &KernelFunctionMetadata{
		Name: name,
	}
}

func KernelFunctionMetadataClone(meta KernelFunctionMetadata) *KernelFunctionMetadata {
	return &KernelFunctionMetadata{
		Name:                 meta.Name,
		PluginName:           meta.PluginName,
		Description:          meta.Description,
		Parameters:           meta.Parameters,
		ReturnParameter:      meta.ReturnParameter,
		AdditionalProperties: meta.AdditionalProperties,
	}
}
